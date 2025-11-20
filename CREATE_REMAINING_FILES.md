# Archivos Restantes para Crear

## Use Cases - Supplier

### 1. create_supplier.go
```go
package supplier

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type CreateSupplierUseCase struct {
	repo ports.SupplierRepository
}

func NewCreateSupplierUseCase(repo ports.SupplierRepository) *CreateSupplierUseCase {
	return &CreateSupplierUseCase{repo: repo}
}

func (uc *CreateSupplierUseCase) Execute(ctx context.Context, supplier *entities.Supplier) error {
	if err := supplier.Validate(); err != nil {
		return err
	}
	supplier.IsActive = true
	return uc.repo.Create(ctx, supplier)
}
```

### 2. get_supplier.go
```go
package supplier

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetSupplierUseCase struct {
	repo ports.SupplierRepository
}

func NewGetSupplierUseCase(repo ports.SupplierRepository) *GetSupplierUseCase {
	return &GetSupplierUseCase{repo: repo}
}

func (uc *GetSupplierUseCase) Execute(ctx context.Context, id uint) (*entities.Supplier, error) {
	return uc.repo.GetByID(ctx, id)
}
```

### 3. list_suppliers.go
```go
package supplier

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type ListSuppliersUseCase struct {
	repo ports.SupplierRepository
}

func NewListSuppliersUseCase(repo ports.SupplierRepository) *ListSuppliersUseCase {
	return &ListSuppliersUseCase{repo: repo}
}

func (uc *ListSuppliersUseCase) Execute(ctx context.Context, filters map[string]interface{}) ([]entities.Supplier, error) {
	return uc.repo.List(ctx, filters)
}
```

### 4. update_supplier.go
```go
package supplier

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type UpdateSupplierUseCase struct {
	repo ports.SupplierRepository
}

func NewUpdateSupplierUseCase(repo ports.SupplierRepository) *UpdateSupplierUseCase {
	return &UpdateSupplierUseCase{repo: repo}
}

func (uc *UpdateSupplierUseCase) Execute(ctx context.Context, supplier *entities.Supplier) error {
	if err := supplier.Validate(); err != nil {
		return err
	}
	return uc.repo.Update(ctx, supplier)
}
```

### 5. delete_supplier.go
```go
package supplier

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type DeleteSupplierUseCase struct {
	repo ports.SupplierRepository
}

func NewDeleteSupplierUseCase(repo ports.SupplierRepository) *DeleteSupplierUseCase {
	return &DeleteSupplierUseCase{repo: repo}
}

func (uc *DeleteSupplierUseCase) Execute(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
```

## Use Cases - Capital Injection

### 1. create_injection.go
```go
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
```

### 2. get_injection.go
```go
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
```

### 3. list_injections.go
```go
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
```

### 4. get_total_capital.go
```go
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
```

## Handlers

Continúa en el siguiente archivo...
