package blizzard

import (
	"context"
	"profile-service/internal/entity"
	"profile-service/pkg/dto"
)

type BlizzardRepository interface {
	GetCharacters(ctx context.Context, blizzAccess string) ([]entity.Character, error)
	GetUserData(ctx context.Context, blizzAccess string) (*dto.UserDTO, error)
}
