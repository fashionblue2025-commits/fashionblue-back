package size

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetSizeUseCase struct {
	sizeRepo ports.SizeRepository
}

func NewGetSizeUseCase(sizeRepo ports.SizeRepository) *GetSizeUseCase {
	return &GetSizeUseCase{
		sizeRepo: sizeRepo,
	}
}

func (uc *GetSizeUseCase) Execute(ctx context.Context, id uint) (*entities.Size, error) {
	return uc.sizeRepo.GetByID(ctx, id)
}
