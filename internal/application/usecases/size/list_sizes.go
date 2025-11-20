package size

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type ListSizesUseCase struct {
	sizeRepo ports.SizeRepository
}

func NewListSizesUseCase(sizeRepo ports.SizeRepository) *ListSizesUseCase {
	return &ListSizesUseCase{
		sizeRepo: sizeRepo,
	}
}

func (uc *ListSizesUseCase) Execute(ctx context.Context, filters map[string]interface{}) ([]*entities.Size, error) {
	return uc.sizeRepo.List(ctx, filters)
}
