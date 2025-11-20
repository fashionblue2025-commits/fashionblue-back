package capital_injection

import (
	"context"
	"time"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type CreateInjectionUseCase struct {
	repo ports.CapitalInjectionRepository
}

func NewCreateInjectionUseCase(repo ports.CapitalInjectionRepository) *CreateInjectionUseCase {
	return &CreateInjectionUseCase{repo: repo}
}

func (uc *CreateInjectionUseCase) Execute(ctx context.Context, injection *entities.CapitalInjection) error {
	if err := injection.Validate(); err != nil {
		return err
	}
	if injection.Date.IsZero() {
		injection.Date = time.Now()
	}
	return uc.repo.Create(ctx, injection)
}
