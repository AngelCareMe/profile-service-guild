package logger

import (
	"profile-service/pkg/config"

	"github.com/sirupsen/logrus"
)

func InitLogger(cfg *config.Config) (*logrus.Logger, error) {
	log := logrus.New()

	logLvl, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		log.SetLevel(logrus.InfoLevel)
		log.Warnf("failed parse logger level: %v", err)
	} else {
		log.SetLevel(logLvl)
	}

	log.SetFormatter(&logrus.JSONFormatter{})

	log.WithFields(logrus.Fields{
		"service": "profile",
	}).Info("logger initialized")

	return log, nil
}
