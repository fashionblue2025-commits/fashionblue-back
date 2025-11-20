package category

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type DeleteCategoryUseCase struct {
	categoryRepo ports.CategoryRepository
}

func NewDeleteCategoryUseCase(categoryRepo ports.CategoryRepository) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{categoryRepo: categoryRepo}
}

func (uc *DeleteCategoryUseCase) Execute(ctx context.Context, id uint) error {
	_, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.categoryRepo.Delete(ctx, id)
}
