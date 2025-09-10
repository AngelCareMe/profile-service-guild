package database

import (
	"context"
	"profile-service/internal/entity"
)

type PostgresRepository interface {
	GetCharacters(ctx context.Context, blizzardID string) ([]entity.Character, error)
	SaveCharacters(ctx context.Context, characters []entity.Character) error
	GetCharacterByName(ctx context.Context, charName string) (*entity.Character, error)
	SetMainCharacter(ctx context.Context, blizzardID, charName string) error
}
