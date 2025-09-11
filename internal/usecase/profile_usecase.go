package usecase

import (
	"context"
	"profile-service/internal/entity"
)

type ProfileUsecase interface {
	GetCharacters(ctx context.Context, blizzardID, accessToken, jwtToken string) ([]entity.Character, error)
	RefreshCharacters(ctx context.Context, blizzardID, accessToken, jwtToken string) error
	SetMain(ctx context.Context, blizzardID, charName string) error
	GetGuildByName(ctx context.Context, nameSlug string) (*entity.Guild, error)
	GetMainCharacterByBlizzardID(ctx context.Context, blizzardID string) (*entity.Character, error)
}
