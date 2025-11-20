package category

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type UpdateCategoryUseCase struct {
	categoryRepo ports.CategoryRepository
}

func NewUpdateCategoryUseCase(categoryRepo ports.CategoryRepository) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{categoryRepo: categoryRepo}
}

func (uc *UpdateCategoryUseCase) Execute(ctx context.Context, category *entities.Category) error {
	return uc.categoryRepo.Update(ctx, category)
}
