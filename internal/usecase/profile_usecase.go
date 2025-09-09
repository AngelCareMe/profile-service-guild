package usecase

import (
	"context"
	"profile-service/internal/entity"
)

type ProfileUsecase interface {
	GetCharacters(ctx context.Context, blizzardID, accessToken, jwtToken string) ([]entity.Character, error)
	RefreshCharacters(ctx context.Context, blizzardID, accessToken, jwtToken string) error
}
