package category

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetCategoryUseCase struct {
	categoryRepo ports.CategoryRepository
}

func NewGetCategoryUseCase(categoryRepo ports.CategoryRepository) *GetCategoryUseCase {
	return &GetCategoryUseCase{categoryRepo: categoryRepo}
}

func (uc *GetCategoryUseCase) Execute(ctx context.Context, id uint) (*entities.Category, error) {
	return uc.categoryRepo.GetByID(ctx, id)
}
