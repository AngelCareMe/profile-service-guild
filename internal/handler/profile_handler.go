package handler

import (
	"net/http"
	"profile-service/internal/adapter/blizzard"
	"profile-service/internal/usecase"
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
	token := strings.TrimPrefix(tokenStr, "Bearer ")

	user, err := h.blizzAd.GetUserData(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid access token"})
		return
	}

	characters, err := h.uc.GetCharacters(c.Request.Context(), user.ID, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed parse characters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"characters": characters})
}

func (h *ProfileHandler) RefreshCharacters(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	token := strings.TrimPrefix(tokenStr, "Bearer ")

	user, err := h.blizzAd.GetUserData(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid access token"})
		return
	}

	if err := h.uc.RefreshCharacters(c.Request.Context(), user.ID, token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed refresh characters"})
		return
	}

	characters, err := h.uc.GetCharacters(c.Request.Context(), user.ID, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed parse characters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"characters": characters})
}
