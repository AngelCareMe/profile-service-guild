package usecase

import (
	"context"
	"fmt"
	"profile-service/internal/adapter/blizzard"
	"profile-service/internal/adapter/database"
	"profile-service/internal/entity"
	"profile-service/pkg/errors"

	"github.com/sirupsen/logrus"
)

type profileUsecase struct {
	dbAd    database.PostgresRepository
	blizzAd blizzard.BlizzardRepository
	log     *logrus.Logger
}

func NewProfileUsecase(
	dbAd database.PostgresRepository,
	blizzAd blizzard.BlizzardRepository,
	log *logrus.Logger,
) *profileUsecase {
	return &profileUsecase{
		dbAd:    dbAd,
		blizzAd: blizzAd,
		log:     log,
	}
}

func (uc *profileUsecase) GetCharacters(ctx context.Context, blizzardID, accessToken string) ([]entity.Character, error) {
	if accessToken == "" || blizzardID == "" {
		uc.log.Warn("id or access token is empty")
		return nil, fmt.Errorf("id or access token is empty")
	}

	var dbChars []entity.Character
	var err error
	if blizzardID != "" {
		dbChars, err = uc.dbAd.GetCharacters(ctx, blizzardID)
		if err != nil {
			uc.log.WithError(err).WithField("blizzard_id", blizzardID).Warn("failed get characters from DB")
		}
	}

	if len(dbChars) > 0 {
		uc.log.Debugf("Found %d characters in DB for blizzard_id %s", len(dbChars), blizzardID)
		return dbChars, nil
	}

	uc.log.Debug("No characters found in DB or blizzardID not provided, fetching from Blizzard API")
	blizzChars, err := uc.blizzAd.GetCharacters(ctx, accessToken)
	if err != nil {
		uc.log.WithError(err).Error("failed fetch characters from Blizzard API")
		return nil, err
	}
	uc.log.Debugf("Fetched %d characters from Blizzard API", len(blizzChars))

	if len(blizzChars) > 0 {
		if err := uc.dbAd.SaveCharacters(ctx, blizzChars); err != nil {
			uc.log.WithError(err).Error("Failed save characters")
			return nil, err
		}
	} else {
		uc.log.Info("BlizzardID is empty in fetched characters")
	}

	return blizzChars, nil
}

func (uc *profileUsecase) RefreshCharacters(ctx context.Context, blizzardID, accessToken string) error {
	if accessToken == "" || blizzardID == "" {
		uc.log.Warn("id or access token is empty")
		return errors.NewAppError("id or access token is empty", nil)
	}

	blizzChars, err := uc.blizzAd.GetCharacters(ctx, accessToken)
	if err != nil {
		uc.log.WithError(err).Error("failed fetch characters from Blizzard API")
		return err
	}

	if err := uc.dbAd.SaveCharacters(ctx, blizzChars); err != nil {
		uc.log.WithError(err).Error("failed save characters")
		return err
	}

	return nil
}
