package category

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type ListCategoriesUseCase struct {
	categoryRepo ports.CategoryRepository
}

func NewListCategoriesUseCase(categoryRepo ports.CategoryRepository) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{categoryRepo: categoryRepo}
}

func (uc *ListCategoriesUseCase) Execute(ctx context.Context, filters map[string]interface{}) ([]entities.Category, error) {
	return uc.categoryRepo.List(ctx, filters)
}
