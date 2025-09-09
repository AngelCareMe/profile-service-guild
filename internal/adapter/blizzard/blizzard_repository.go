package blizzard

import (
	"context"
	"profile-service/internal/entity"
	"profile-service/pkg/dto"
)

type BlizzardRepository interface {
	GetCharacters(ctx context.Context, blizzAccess, jwtToken string) ([]entity.Character, error)
	GetUserData(ctx context.Context, jwtToken string) (*dto.UserDTO, error)
	GetBlizzardAccessToken(ctx context.Context, jwtToken string) (string, error)
}
