# Guía de Refactorización - Fashion Blue

## 🎯 Objetivo de la Refactorización

Separar completamente las **entidades de dominio** de los **modelos de persistencia** y reorganizar la estructura de carpetas para mejor mantenibilidad.

## 📐 Nueva Arquitectura

### Antes vs Después

#### **Antes:**
```
internal/domain/
├── user.go              # Entidad con tags GORM
├── product.go           # Entidad con tags GORM
└── ...
```

#### **Después:**
```
internal/
├── domain/
│   ├── entities/        # Entidades puras (sin GORM)
│   │   ├── user.go
│   │   ├── product.go
│   │   └── ...
│   └── ports/           # Interfaces organizadas
│       ├── user_repository.go
│       ├── product_repository.go
│       └── services.go
├── application/
│   └── usecases/        # Casos de uso por entidad
│       ├── user/
│       │   ├── create_user.go
│       │   ├── get_user.go
│       │   └── ...
│       ├── product/
│       │   ├── create_product.go
│       │   └── ...
│       └── ...
└── adapters/
    ├── persistence/
    │   ├── models/      # Modelos GORM
    │   │   ├── user_model.go
    │   │   ├── product_model.go
    │   │   └── ...
    │   └── repositories/  # Implementaciones
    │       ├── user/
    │       │   └── user_repository.go
    │       ├── product/
    │       │   └── product_repository.go
    │       └── ...
    └── http/
        └── handlers/
            ├── user/
            │   └── user_handler.go
            ├── product/
            │   └── product_handler.go
            └── ...
```

## ✅ Cambios Completados

### 1. Entidades de Dominio Puras
- ✅ Creadas en `internal/domain/entities/`
- ✅ Sin dependencias de GORM
- ✅ Solo lógica de negocio
- ✅ Archivos creados:
  - `user.go`
  - `customer.go`
  - `product.go`
  - `category.go`
  - `sale.go`
  - `supplier.go`
  - `capital_injection.go`

### 2. Modelos de Persistencia
- ✅ Creados en `internal/adapters/persistence/models/`
- ✅ Con tags GORM
- ✅ Métodos `ToEntity()` y `FromEntity()` para conversión
- ✅ Archivos creados:
  - `user_model.go`
  - `customer_model.go`
  - `product_model.go`
  - `sale_model.go`
  - `supplier_model.go`
  - `capital_injection_model.go`

### 3. Puertos (Interfaces) Reorganizados
- ✅ Separados por dominio en `internal/domain/ports/`
- ✅ Usan entidades puras
- ✅ Archivos creados:
  - `user_repository.go`
  - `product_repository.go`
  - `customer_repository.go`
  - `sale_repository.go`
  - `supplier_repository.go`
  - `capital_injection_repository.go`
  - `services.go`

## 🚧 Pendiente de Refactorizar

### 1. Reorganizar Casos de Uso
Mover de:
```
internal/application/services/user_service.go
```

A:
```
internal/application/usecases/user/
├── create_user.go
├── get_user.go
├── list_users.go
├── update_user.go
├── delete_user.go
└── change_password.go
```

**Estructura de cada caso de uso:**
```go
package user

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type CreateUserUseCase struct {
	userRepo ports.UserRepository
}

func NewCreateUserUseCase(userRepo ports.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: userRepo}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, user *entities.User, password string) error {
	// Lógica del caso de uso
	if err := user.HashPassword(password); err != nil {
		return err
	}
	return uc.userRepo.Create(ctx, user)
}
```

### 2. Reorganizar Repositorios
Mover de:
```
internal/adapters/postgres/user_repository.go
```

A:
```
internal/adapters/persistence/repositories/user/
└── user_repository.go
```

**Actualizar para usar modelos:**
```go
package user

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	model := &models.UserModel{}
	model.FromEntity(user)
	
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	
	*user = *model.ToEntity()
	return nil
}
```

### 3. Reorganizar Handlers
Mover de:
```
internal/adapters/http/handlers/user_handler.go
```

A:
```
internal/adapters/http/handlers/user/
└── user_handler.go
```

**Actualizar para usar casos de uso:**
```go
package user

import (
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/user"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	createUserUC *user.CreateUserUseCase
	getUserUC    *user.GetUserUseCase
	// ... otros casos de uso
}

func NewUserHandler(
	createUserUC *user.CreateUserUseCase,
	getUserUC *user.GetUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUC: createUserUC,
		getUserUC:    getUserUC,
	}
}
```

### 4. Actualizar Database Package
Actualizar `pkg/database/postgres.go` para usar modelos:

```go
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.UserModel{},
		&models.CapitalInjectionModel{},
		&models.CategoryModel{},
		&models.ProductModel{},
		&models.CustomerModel{},
		&models.CustomerTransactionModel{},
		&models.SaleModel{},
		&models.SaleItemModel{},
		&models.SupplierModel{},
		&models.PurchaseModel{},
		&models.PurchaseItemModel{},
	)
}
```

### 5. Actualizar Scripts
Actualizar `scripts/seed.go` para usar modelos y entidades correctamente.

### 6. Actualizar Main
Actualizar `cmd/api/main.go` para instanciar casos de uso individuales.

## 📝 Pasos para Completar la Refactorización

### Paso 1: Reorganizar User (Ejemplo Completo)
1. Crear casos de uso en `internal/application/usecases/user/`
2. Crear repositorio en `internal/adapters/persistence/repositories/user/`
3. Crear handler en `internal/adapters/http/handlers/user/`
4. Actualizar rutas y main.go

### Paso 2: Aplicar el mismo patrón a las demás entidades
- Product
- Customer
- Sale
- Supplier
- Purchase
- CapitalInjection
- Category

### Paso 3: Limpiar archivos antiguos
Una vez verificado que todo funciona:
```bash
rm -rf internal/domain/user.go
rm -rf internal/domain/product.go
# ... etc

rm -rf internal/adapters/postgres/
rm -rf internal/application/services/
rm -rf internal/adapters/http/handlers/*.go
```

### Paso 4: Actualizar imports en toda la aplicación
Buscar y reemplazar:
- `internal/domain` → `internal/domain/entities`
- `internal/ports` → `internal/domain/ports`

## 🎨 Beneficios de la Nueva Arquitectura

### 1. **Separación de Responsabilidades**
- Dominio no conoce la persistencia
- Fácil cambiar de ORM o base de datos
- Testeable sin dependencias externas

### 2. **Organización Clara**
- Un archivo por caso de uso
- Fácil encontrar funcionalidad específica
- Mejor para equipos grandes

### 3. **Mantenibilidad**
- Cambios localizados
- Menos acoplamiento
- Más fácil de extender

### 4. **Testing**
- Casos de uso independientes
- Mock de repositorios simple
- Tests unitarios más claros

## 🔄 Flujo de Datos

```
HTTP Request
    ↓
Handler (adapters/http/handlers/user/)
    ↓
Use Case (application/usecases/user/)
    ↓
Repository Interface (domain/ports/)
    ↓
Repository Implementation (adapters/persistence/repositories/user/)
    ↓
Model (adapters/persistence/models/)
    ↓
GORM → PostgreSQL
```

## 📚 Ejemplo Completo: CreateUser

### 1. Entidad (domain/entities/user.go)
```go
type User struct {
	ID        uint
	Email     string
	FirstName string
	// ... sin tags GORM
}
```

### 2. Modelo (adapters/persistence/models/user_model.go)
```go
type UserModel struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex"`
	// ... con tags GORM
}

func (m *UserModel) ToEntity() *entities.User { ... }
func (m *UserModel) FromEntity(user *entities.User) { ... }
```

### 3. Port (domain/ports/user_repository.go)
```go
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
}
```

### 4. Caso de Uso (application/usecases/user/create_user.go)
```go
type CreateUserUseCase struct {
	userRepo ports.UserRepository
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, user *entities.User) error {
	return uc.userRepo.Create(ctx, user)
}
```

### 5. Repositorio (adapters/persistence/repositories/user/user_repository.go)
```go
func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	model := &models.UserModel{}
	model.FromEntity(user)
	return r.db.Create(model).Error
}
```

### 6. Handler (adapters/http/handlers/user/user_handler.go)
```go
func (h *UserHandler) Create(c echo.Context) error {
	var user entities.User
	c.Bind(&user)
	return h.createUserUC.Execute(c.Request().Context(), &user)
}
```

## ⚠️ Notas Importantes

1. **No borrar archivos antiguos hasta verificar** que todo funciona
2. **Hacer la migración por partes** (entidad por entidad)
3. **Mantener tests actualizados** durante la refactorización
4. **Documentar cambios** en cada PR

## 🚀 Próximos Pasos Inmediatos

1. ✅ Entidades puras creadas
2. ✅ Modelos de persistencia creados
3. ✅ Puertos reorganizados
4. ⏳ Crear ejemplo completo de User refactorizado
5. ⏳ Aplicar patrón a las demás entidades
6. ⏳ Actualizar tests
7. ⏳ Limpiar código antiguo

¿Quieres que continúe con la implementación completa de un ejemplo (User) para que veas el patrón completo?
