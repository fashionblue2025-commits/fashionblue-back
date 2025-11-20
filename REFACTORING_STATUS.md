# Estado de la Refactorización - Fashion Blue

## ✅ Completado

### 1. Separación de Entidades y Modelos

#### **Entidades de Dominio** (`internal/domain/entities/`)
Entidades puras sin dependencias de GORM:
- ✅ `user.go` - Usuario con lógica de negocio
- ✅ `customer.go` - Cliente y transacciones
- ✅ `product.go` - Producto con cálculos de ganancia
- ✅ `category.go` - Categoría de productos
- ✅ `sale.go` - Venta e ítems de venta
- ✅ `supplier.go` - Proveedor, compra e ítems
- ✅ `capital_injection.go` - Inyección de capital

#### **Modelos de Persistencia** (`internal/adapters/persistence/models/`)
Modelos con tags GORM y conversión ToEntity/FromEntity:
- ✅ `user_model.go`
- ✅ `customer_model.go`
- ✅ `product_model.go`
- ✅ `sale_model.go`
- ✅ `supplier_model.go`
- ✅ `capital_injection_model.go`

### 2. Reorganización de Puertos

#### **Puertos por Dominio** (`internal/domain/ports/`)
- ✅ `user_repository.go`
- ✅ `product_repository.go`
- ✅ `customer_repository.go`
- ✅ `sale_repository.go`
- ✅ `supplier_repository.go`
- ✅ `capital_injection_repository.go`
- ✅ `services.go` - Interfaces de servicios

### 3. Ejemplo Completo: User

#### **Casos de Uso** (`internal/application/usecases/user/`)
- ✅ `create_user.go` - Crear usuario con validaciones
- ✅ `get_user.go` - Obtener usuario por ID
- ✅ `list_users.go` - Listar usuarios con filtros
- ✅ `update_user.go` - Actualizar usuario
- ✅ `delete_user.go` - Eliminar usuario
- ✅ `change_password.go` - Cambiar contraseña

#### **Repositorio** (`internal/adapters/persistence/repositories/user/`)
- ✅ `user_repository.go` - Implementación con conversión de modelos

#### **Handler** (`internal/adapters/http/handlers/user/`)
- ✅ `user_handler.go` - Handler HTTP usando casos de uso

## 📊 Estructura Actual

```
fashion-blue/
├── internal/
│   ├── domain/
│   │   ├── entities/              ✅ NUEVO - Entidades puras
│   │   │   ├── user.go
│   │   │   ├── customer.go
│   │   │   ├── product.go
│   │   │   ├── category.go
│   │   │   ├── sale.go
│   │   │   ├── supplier.go
│   │   │   └── capital_injection.go
│   │   ├── ports/                 ✅ REORGANIZADO
│   │   │   ├── user_repository.go
│   │   │   ├── product_repository.go
│   │   │   ├── customer_repository.go
│   │   │   ├── sale_repository.go
│   │   │   ├── supplier_repository.go
│   │   │   ├── capital_injection_repository.go
│   │   │   └── services.go
│   │   ├── user.go                ⚠️ ANTIGUO - Mantener por ahora
│   │   ├── product.go             ⚠️ ANTIGUO - Mantener por ahora
│   │   └── ...
│   ├── application/
│   │   ├── usecases/              ✅ NUEVO
│   │   │   └── user/              ✅ COMPLETO
│   │   │       ├── create_user.go
│   │   │       ├── get_user.go
│   │   │       ├── list_users.go
│   │   │       ├── update_user.go
│   │   │       ├── delete_user.go
│   │   │       └── change_password.go
│   │   └── services/              ⚠️ ANTIGUO - Mantener por ahora
│   │       ├── user_service.go
│   │       └── ...
│   └── adapters/
│       ├── persistence/
│       │   ├── models/            ✅ NUEVO - Modelos GORM
│       │   │   ├── user_model.go
│       │   │   ├── customer_model.go
│       │   │   ├── product_model.go
│       │   │   ├── sale_model.go
│       │   │   ├── supplier_model.go
│       │   │   └── capital_injection_model.go
│       │   └── repositories/      ✅ NUEVO
│       │       └── user/          ✅ COMPLETO
│       │           └── user_repository.go
│       ├── postgres/              ⚠️ ANTIGUO - Mantener por ahora
│       │   ├── user_repository.go
│       │   └── ...
│       └── http/
│           └── handlers/
│               ├── user/          ✅ NUEVO - COMPLETO
│               │   └── user_handler.go
│               ├── user_handler.go ⚠️ ANTIGUO - Mantener por ahora
│               └── ...
```

## 🎯 Patrón Implementado (User)

### Flujo de Datos

```
HTTP Request
    ↓
UserHandler (adapters/http/handlers/user/)
    ↓
CreateUserUseCase (application/usecases/user/)
    ↓
UserRepository Interface (domain/ports/)
    ↓
UserRepository Implementation (adapters/persistence/repositories/user/)
    ↓
UserModel (adapters/persistence/models/)
    ↓
GORM → PostgreSQL
```

### Ejemplo de Conversión

```go
// 1. HTTP Request llega al Handler
handler.Create(c echo.Context)

// 2. Handler crea entidad de dominio
user := &entities.User{
    Email: "user@example.com",
    // ... sin tags GORM
}

// 3. Llama al caso de uso
createUserUC.Execute(ctx, user, password)

// 4. Caso de uso valida y llama al repositorio
userRepo.Create(ctx, user)

// 5. Repositorio convierte a modelo
model := &models.UserModel{}
model.FromEntity(user)  // Convierte entidad → modelo

// 6. GORM persiste el modelo
db.Create(model)

// 7. Repositorio convierte de vuelta
*user = *model.ToEntity()  // Convierte modelo → entidad

// 8. Handler devuelve la entidad
return response.Created(c, "User created", user)
```

## 📝 Próximos Pasos

### Fase 1: Completar Refactorización de Entidades Restantes

#### Product
- [ ] Crear casos de uso en `usecases/product/`
- [ ] Crear repositorio en `repositories/product/`
- [ ] Crear handler en `handlers/product/`

#### Customer
- [ ] Crear casos de uso en `usecases/customer/`
- [ ] Crear repositorio en `repositories/customer/`
- [ ] Crear handler en `handlers/customer/`

#### Sale
- [ ] Crear casos de uso en `usecases/sale/`
- [ ] Crear repositorio en `repositories/sale/`
- [ ] Crear handler en `handlers/sale/`

#### Supplier
- [ ] Crear casos de uso en `usecases/supplier/`
- [ ] Crear repositorio en `repositories/supplier/`
- [ ] Crear handler en `handlers/supplier/`

#### Purchase
- [ ] Crear casos de uso en `usecases/purchase/`
- [ ] Crear repositorio en `repositories/purchase/`
- [ ] Crear handler en `handlers/purchase/`

#### CapitalInjection
- [ ] Crear casos de uso en `usecases/capital_injection/`
- [ ] Crear repositorio en `repositories/capital_injection/`
- [ ] Crear handler en `handlers/capital_injection/`

#### Category
- [ ] Crear casos de uso en `usecases/category/`
- [ ] Crear repositorio en `repositories/category/`
- [ ] Crear handler en `handlers/category/`

### Fase 2: Actualizar Infraestructura

#### Auth Service
- [ ] Refactorizar para usar entidades puras
- [ ] Mover a `usecases/auth/`

#### Database
- [ ] Actualizar `pkg/database/postgres.go` para usar modelos
- [ ] Actualizar AutoMigrate

#### Seed Script
- [ ] Actualizar `scripts/seed.go` para usar modelos y entidades

#### Main
- [ ] Actualizar `cmd/api/main.go` para instanciar casos de uso
- [ ] Actualizar rutas

### Fase 3: Limpieza

Una vez todo funcione:
- [ ] Eliminar `internal/domain/*.go` (archivos antiguos)
- [ ] Eliminar `internal/adapters/postgres/`
- [ ] Eliminar `internal/application/services/`
- [ ] Eliminar `internal/adapters/http/handlers/*.go` (archivos antiguos)
- [ ] Eliminar `internal/ports/repositories.go` y `services.go` antiguos

### Fase 4: Testing

- [ ] Crear tests unitarios para casos de uso
- [ ] Crear tests de integración para repositorios
- [ ] Crear tests E2E para handlers

## 🔍 Cómo Aplicar el Patrón a Otras Entidades

### 1. Casos de Uso

Crear un archivo por operación en `internal/application/usecases/{entity}/`:

```go
// create_{entity}.go
type Create{Entity}UseCase struct {
    repo ports.{Entity}Repository
}

func (uc *Create{Entity}UseCase) Execute(ctx context.Context, entity *entities.{Entity}) error {
    // Validaciones de negocio
    // Llamar al repositorio
    return uc.repo.Create(ctx, entity)
}
```

### 2. Repositorio

Crear en `internal/adapters/persistence/repositories/{entity}/`:

```go
// {entity}_repository.go
type {entity}Repository struct {
    db *gorm.DB
}

func (r *{entity}Repository) Create(ctx context.Context, entity *entities.{Entity}) error {
    model := &models.{Entity}Model{}
    model.FromEntity(entity)
    
    if err := r.db.Create(model).Error; err != nil {
        return err
    }
    
    *entity = *model.ToEntity()
    return nil
}
```

### 3. Handler

Crear en `internal/adapters/http/handlers/{entity}/`:

```go
// {entity}_handler.go
type {Entity}Handler struct {
    create{Entity}UC *{entity}.Create{Entity}UseCase
    // ... otros casos de uso
}

func (h *{Entity}Handler) Create(c echo.Context) error {
    var entity entities.{Entity}
    c.Bind(&entity)
    
    err := h.create{Entity}UC.Execute(c.Request().Context(), &entity)
    if err != nil {
        return response.BadRequest(c, "Failed", err)
    }
    
    return response.Created(c, "Created", entity)
}
```

## 💡 Beneficios Observados

### 1. Separación Clara
- Dominio no conoce GORM ✅
- Fácil cambiar de ORM ✅
- Testeable sin DB ✅

### 2. Organización
- Un archivo por caso de uso ✅
- Fácil encontrar código ✅
- Mejor para equipos ✅

### 3. Mantenibilidad
- Cambios localizados ✅
- Menos acoplamiento ✅
- Más fácil extender ✅

## 📚 Documentación Relacionada

- `REFACTORING_GUIDE.md` - Guía completa de refactorización
- `README.md` - Documentación principal
- `SETUP.md` - Guía de instalación

## 🚀 Comandos Útiles

```bash
# Ver estructura de archivos nuevos
find internal/domain/entities -type f
find internal/adapters/persistence/models -type f
find internal/application/usecases -type f

# Buscar referencias a archivos antiguos
grep -r "internal/domain" --include="*.go" | grep -v "entities" | grep -v "ports"

# Ejecutar tests (cuando estén creados)
go test ./internal/application/usecases/user/...
go test ./internal/adapters/persistence/repositories/user/...
```

## ⚠️ Notas Importantes

1. **No borrar archivos antiguos** hasta que todo esté migrado y probado
2. **Mantener ambas implementaciones** funcionando en paralelo
3. **Migrar entidad por entidad** para minimizar riesgos
4. **Actualizar tests** a medida que se refactoriza
5. **Documentar cambios** en cada commit

---

**Última actualización:** Fase 1 completada para User
**Siguiente paso:** Aplicar patrón a Product
