#!/bin/bash

# Script para restaurar Suppliers y Capital Injections

echo "🔧 Creando archivos de Suppliers y Capital Injections..."

# Use Cases - Supplier
cat > internal/application/usecases/supplier/create_supplier.go << 'EOF'
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
EOF

cat > internal/application/usecases/supplier/get_supplier.go << 'EOF'
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
EOF

cat > internal/application/usecases/supplier/list_suppliers.go << 'EOF'
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
EOF

cat > internal/application/usecases/supplier/update_supplier.go << 'EOF'
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
EOF

cat > internal/application/usecases/supplier/delete_supplier.go << 'EOF'
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
EOF

# Use Cases - Capital Injection
cat > internal/application/usecases/capital_injection/create_injection.go << 'EOF'
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
EOF

cat > internal/application/usecases/capital_injection/get_injection.go << 'EOF'
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
EOF

cat > internal/application/usecases/capital_injection/list_injections.go << 'EOF'
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
EOF

cat > internal/application/usecases/capital_injection/get_total_capital.go << 'EOF'
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
EOF

echo "✅ Archivos de use cases creados"
echo "📝 Ahora ejecuta: chmod +x restore_suppliers_capital.sh && ./restore_suppliers_capital.sh"
