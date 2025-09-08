package database

import (
	"context"
	"profile-service/internal/entity"
)

type PostgresRepository interface {
	GetCharacters(ctx context.Context, blizzardID string) ([]entity.Character, error)
	SaveCharacters(ctx context.Context, characters []entity.Character) error
}
