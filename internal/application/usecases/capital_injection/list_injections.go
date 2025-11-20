package capital_injection

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type ListInjectionsUseCase struct {
	repo ports.CapitalInjectionRepository
}

func NewListInjectionsUseCase(repo ports.CapitalInjectionRepository) *ListInjectionsUseCase {
	return &ListInjectionsUseCase{repo: repo}
}

func (uc *ListInjectionsUseCase) Execute(ctx context.Context, filters map[string]interface{}) ([]entities.CapitalInjection, error) {
	return uc.repo.List(ctx, filters)
}
