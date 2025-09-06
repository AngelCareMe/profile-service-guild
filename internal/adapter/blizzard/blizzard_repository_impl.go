package blizzard

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type blizzardRepository struct {
	client *http.Client
	log    *logrus.Logger
}

func NewBlizzardRepository(client *http.Client, log *logrus.Logger) *blizzardRepository{
	return &blizzardRepository{
		client: client,
		log: log,
	}
}

func(br *blizzardRepository) getMythicScore
