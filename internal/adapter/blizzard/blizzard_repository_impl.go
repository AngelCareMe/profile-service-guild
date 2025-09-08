package blizzard

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"profile-service/internal/entity"
	"profile-service/pkg/dto"
	"profile-service/pkg/errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type blizzardRepository struct {
	client  *http.Client
	log     *logrus.Logger
	limiter *rate.Limiter
}

func NewBlizzardRepository(log *logrus.Logger) *blizzardRepository {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			MaxConnsPerHost:    5,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: false,
		},
	}
	return &blizzardRepository{
		client:  client,
		log:     log,
		limiter: rate.NewLimiter(rate.Every(time.Second/50), 50),
	}
}

func (br *blizzardRepository) GetUserData(ctx context.Context, blizzAccess string) (*dto.UserDTO, error) {
	if blizzAccess == "" {
		br.log.Error("access header is missing")
		return nil, errors.NewAppError("access token is empty", nil)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/auth/user", nil)
	if err != nil {
		br.log.WithError(err).Errorf("failed create get user request: %v ", err)
		return nil, errors.NewAppError("failed create get user request", err)
	}
	req.Header.Set("Authorization", "Bearer "+blizzAccess)

	resp, err := br.client.Do(req)
	if err != nil {
		br.log.WithError(err).WithField("req", req).Error("failed get user response")
		return nil, errors.NewAppError("failed get user response", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		br.log.WithFields(logrus.Fields{
			"status": resp.StatusCode,
			"body":   string(body),
		}).Warn("bad response from auth service")
		return nil, errors.NewAppError("bad response", err)
	}

	var userDTO dto.UserDTO
	if err := json.NewDecoder(resp.Body).Decode(&userDTO); err != nil {
		br.log.WithError(err).Error("failed to decode user response")
		return nil, errors.NewAppError("failed to user profile", err)
	}

	return &userDTO, nil
}

func (br *blizzardRepository) GetCharacters(ctx context.Context, blizzAccess string) ([]entity.Character, error) {
	if blizzAccess == "" {
		br.log.Error("access header is missing")
		return nil, errors.NewAppError("access token is empty", nil)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://eu.api.blizzard.com/profile/user/wow?namespace=profile-eu&locale=ru_RU", nil)
	if err != nil {
		br.log.WithError(err).Errorf("failed create characters request: %v ", err)
		return nil, errors.NewAppError("failed create characters request", err)
	}
	req.Header.Set("Authorization", "Bearer "+blizzAccess)

	resp, err := br.client.Do(req)
	if err != nil {
		br.log.WithError(err).WithField("req", req).Error("failed get characters response")
		return nil, errors.NewAppError("failed get characters response", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		br.log.WithFields(logrus.Fields{
			"status": resp.StatusCode,
			"body":   string(body),
		}).Warn("bad response from API")
		return nil, errors.NewAppError("bad response", err)
	}

	var profile dto.BlizzardProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		br.log.WithError(err).Error("failed to decode profile response")
		return nil, errors.NewAppError("failed to decode profile", err)
	}

	var characters []entity.Character
	var mu sync.Mutex
	var wg sync.WaitGroup

	user, err := br.GetUserData(ctx, blizzAccess)
	if err != nil {
		br.log.WithError(err).Error("failed parse user data")
		return nil, err
	}

	for _, acc := range profile.WowAccounts {
		for _, char := range acc.Characters {
			wg.Add(1)
			go func(c dto.CharacterSummary) {
				defer wg.Done()

				mythScore, err := br.getMythicScore(ctx, blizzAccess, char.Realm.Slug, char.Name)
				if err != nil {
					br.log.WithError(err).WithFields(logrus.Fields{
						"character": c.Name,
						"realm":     c.Realm.Slug,
					}).Error("failed parse mythic score")
					mythScore = 0
				}

				details, err := br.getCharacterDetails(ctx, blizzAccess, char.Realm.Slug, char.Name)
				if err != nil {
					br.log.WithError(err).WithFields(logrus.Fields{
						"character": c.Name,
						"realm":     c.Realm.Slug,
					}).Error("failed parse character details")
					return
				}

				character := entity.Character{
					BlizzardID:  user.ID,
					Battletag:   user.Battletag,
					Name:        char.Name,
					Race:        char.PlayableRace.Name,
					Realm:       char.Realm.Name,
					Faction:     char.Faction.Name,
					Class:       char.PlayableClass.Name,
					Spec:        details.Spec.Name,
					Lvl:         char.Level,
					Ilvl:        details.Ilvl,
					Guild:       details.Guild.Name,
					MythicScore: mythScore,
					IsMain:      false,
				}

				mu.Lock()
				characters = append(characters, character)
				mu.Unlock()
			}(char)
		}
	}
	wg.Wait()

	br.log.Info("Characters parsed succeeded")
	return characters, nil
}

func (br *blizzardRepository) getCharacterDetails(ctx context.Context, blizzAccess, realm, charName string) (*dto.CharacterDetailsResponse, error) {
	if blizzAccess == "" || realm == "" || charName == "" {
		br.log.Error("access header/realm/character name is missing")
		return nil, fmt.Errorf("token/realm/character name is empty")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://eu.api.blizzard.com/profile/wow/character/%s/%s?namespace=profile-eu&locale=ru_RU", realm, charName), nil)
	if err != nil {
		br.log.WithError(err).Errorf("failed create character details request for character %s: %v ", charName, err)
		return nil, errors.NewAppError("failed create character details request", err)
	}
	req.Header.Set("Authorization", "Bearer "+blizzAccess)

	if err := br.limiter.Wait(ctx); err != nil {
		return nil, errors.NewAppError("rate limit exeeded: %w", err)
	}

	resp, err := br.client.Do(req)
	if err != nil {
		br.log.WithError(err).WithFields(logrus.Fields{
			"req":       req,
			"character": charName,
			"realm":     realm,
		}).Error("failed get character details response by api")
		return nil, errors.NewAppError("failed get character details response by api", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		br.log.WithField("status", resp.StatusCode).Warn("bad request")
		return nil, fmt.Errorf("bad request")
	}

	var details dto.CharacterDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		br.log.WithError(err).WithField("body", resp.Body).Warn("failed to decode details")
		return nil, errors.NewAppError("failed to decode details", err)
	}

	if details.Spec.Name == "" || details.Ilvl == 0 {
		br.log.Warn("incomplete details response")
	}

	return &details, nil
}

func (br *blizzardRepository) getMythicScore(ctx context.Context, blizzAccess, realm, charName string) (float64, error) {
	if blizzAccess == "" || realm == "" || charName == "" {
		br.log.Error("access header/realm/character name is missing")
		return 0, fmt.Errorf("token/realm/character name is empty")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://eu.api.blizzard.com/profile/wow/character/%s/%s/mythic-keystone-profile?namespace=profile-eu&locale=ru_RU", realm, charName), nil)
	if err != nil {
		br.log.WithError(err).Errorf("failed create mythic score request for character %s: %v ", charName, err)
		return 0, errors.NewAppError("failed create mythic score request", err)
	}
	req.Header.Set("Authorization", "Bearer "+blizzAccess)

	if err := br.limiter.Wait(ctx); err != nil {
		return 0, errors.NewAppError("rate limit exeeded: %w", err)
	}

	resp, err := br.client.Do(req)
	if err != nil {
		br.log.WithError(err).WithFields(logrus.Fields{
			"req":       req,
			"character": charName,
			"realm":     realm,
		}).Error("failed get mythic score response by api")
		return 0, errors.NewAppError("failed get mythic score response by api", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, nil
	}

	if resp.StatusCode != http.StatusOK {
		br.log.WithField("status", resp.StatusCode).Warn("bad request")
		return 0, fmt.Errorf("bad request")
	}

	var mythScoreDto dto.MythScoreDto
	if err := json.NewDecoder(resp.Body).Decode(&mythScoreDto); err != nil {
		br.log.WithError(err).Warn("failed to decode mythic score")
		return 0, errors.NewAppError("failed to decode mythic score", err)
	}

	return mythScoreDto.Current.Rating, nil
}
