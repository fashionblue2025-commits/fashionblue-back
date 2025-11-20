package capital_injection

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetTotalCapitalUseCase struct {
	repo ports.CapitalInjectionRepository
}

func NewGetTotalCapitalUseCase(repo ports.CapitalInjectionRepository) *GetTotalCapitalUseCase {
	return &GetTotalCapitalUseCase{repo: repo}
}

func (uc *GetTotalCapitalUseCase) Execute(ctx context.Context) (float64, error) {
	return uc.repo.GetTotal(ctx)
}
