package size

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetSizesByTypeUseCase struct {
	sizeRepo ports.SizeRepository
}

func NewGetSizesByTypeUseCase(sizeRepo ports.SizeRepository) *GetSizesByTypeUseCase {
	return &GetSizesByTypeUseCase{
		sizeRepo: sizeRepo,
	}
}

func (uc *GetSizesByTypeUseCase) Execute(ctx context.Context, sizeType entities.SizeType) ([]*entities.Size, error) {
	return uc.sizeRepo.GetByType(ctx, sizeType)
}
