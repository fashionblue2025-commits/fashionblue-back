package capital_injection

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetInjectionUseCase struct {
	repo ports.CapitalInjectionRepository
}

func NewGetInjectionUseCase(repo ports.CapitalInjectionRepository) *GetInjectionUseCase {
	return &GetInjectionUseCase{repo: repo}
}

func (uc *GetInjectionUseCase) Execute(ctx context.Context, id uint) (*entities.CapitalInjection, error) {
	return uc.repo.GetByID(ctx, id)
}
