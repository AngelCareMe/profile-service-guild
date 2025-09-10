package handler

import (
	"math"
	"net/http"
	"profile-service/internal/adapter/blizzard"
	"profile-service/internal/usecase"
	"profile-service/pkg/dto"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProfileHandler struct {
	blizzAd blizzard.BlizzardRepository
	uc      usecase.ProfileUsecase
	log     *logrus.Logger
}

func NewProfileHandler(
	blizzAd blizzard.BlizzardRepository,
	uc usecase.ProfileUsecase,
	log *logrus.Logger,
) *ProfileHandler {
	return &ProfileHandler{
		blizzAd: blizzAd,
		uc:      uc,
		log:     log,
	}
}

func (h *ProfileHandler) GetCharacters(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		h.log.Error("Auth header is missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth header"})
		return
	}
	token := strings.TrimPrefix(tokenStr, "Bearer ")
	if token == "" {
		h.log.Error("Invalid Bearer token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Bearer token"})
		return
	}

	tokenAccess, err := h.blizzAd.GetBlizzardAccessToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid access token"})
		return
	}

	user, err := h.blizzAd.GetUserData(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid access token"})
		return
	}

	characters, err := h.uc.GetCharacters(c.Request.Context(), user.ID, tokenAccess, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed parse characters"})
		return
	}

	characterResponses := make([]dto.CharacterResponse, len(characters))
	for i, char := range characters {
		characterResponses[i] = dto.CharacterResponse{
			Name:        char.Name,
			Realm:       char.Realm,
			Race:        char.Race,
			Faction:     char.Faction,
			Class:       char.Class,
			Spec:        char.Spec,
			Lvl:         char.Lvl,
			Ilvl:        char.Ilvl,
			Guild:       char.Guild,
			MythicScore: math.Round(char.MythicScore*100) / 100,
			IsMain:      char.IsMain,
		}
	}

	profileResponse := dto.ProfileResponse{
		BlizzardID: user.ID,
		Battletag:  user.Battletag,
		Characters: characterResponses,
	}

	c.JSON(http.StatusOK, profileResponse)
}

func (h ProfileHandler) SetMainCharacter(c *gin.Context) {
	charName := c.Query("character")
	if charName == "" {
		h.log.Error("Character name is missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing char name"})
		return
	}

	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		h.log.Error("Auth header is missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth header"})
		return
	}

	token := strings.TrimPrefix(jwtToken, "Bearer ")
	if token == "" {
		h.log.Error("Invalid Bearer token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Bearer token"})
		return
	}

	user, err := h.blizzAd.GetUserData(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid access token"})
		return
	}

	if err := h.uc.SetMain(c.Request.Context(), user.ID, charName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed set main character"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Main character set succeed",
		"character": charName,
	})
}

func (h *ProfileHandler) RefreshCharacters(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		h.log.Error("Auth header is missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth header"})
		return
	}
	token := strings.TrimPrefix(tokenStr, "Bearer ")
	if token == "" {
		h.log.Error("Invalid Bearer token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Bearer token"})
		return
	}

	user, err := h.blizzAd.GetUserData(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid access token"})
		return
	}

	tokenAccess, err := h.blizzAd.GetBlizzardAccessToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid access token"})
		return
	}

	if err := h.uc.RefreshCharacters(c.Request.Context(), user.ID, tokenAccess, token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed refresh characters"})
		return
	}

	characters, err := h.uc.GetCharacters(c.Request.Context(), user.ID, tokenAccess, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed parse characters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"characters": characters})
}
