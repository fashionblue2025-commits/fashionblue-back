package category

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type CreateCategoryUseCase struct {
	categoryRepo ports.CategoryRepository
}

func NewCreateCategoryUseCase(categoryRepo ports.CategoryRepository) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{categoryRepo: categoryRepo}
}

func (uc *CreateCategoryUseCase) Execute(ctx context.Context, category *entities.Category) error {
	return uc.categoryRepo.Create(ctx, category)
}
