package handler

import (
	"profile-service/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(
	router *gin.Engine,
	h *ProfileHandler,
	cfg *config.Config,
	log *logrus.Logger,
) {
	router.Use(gin.Recovery())

	profile := router.Group("/profile")
	profile.GET("/refresh", h.RefreshCharacters)
	profile.GET("/characters", h.GetCharacters)
}
