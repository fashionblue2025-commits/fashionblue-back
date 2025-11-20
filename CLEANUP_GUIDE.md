# 🧹 Guía de Limpieza - Archivos Antiguos

## 📋 Resumen

Después de la refactorización a arquitectura hexagonal, hay archivos antiguos que ya no se usan y deben eliminarse.

---

## 🗂️ Archivos a Eliminar

### 1. **Handlers Antiguos** (en raíz de handlers/)

Estos handlers están **duplicados**. Los nuevos están en subcarpetas.

| ❌ Archivo Antiguo | ✅ Archivo Nuevo |
|-------------------|------------------|
| `internal/adapters/http/handlers/auth_handler.go` | `internal/adapters/http/handlers/auth/auth_handler.go` |
| `internal/adapters/http/handlers/user_handler.go` | `internal/adapters/http/handlers/user/user_handler.go` |
| `internal/adapters/http/handlers/capital_injection_handler.go` | `internal/adapters/http/handlers/capital_injection/capital_injection_handler.go` |
| `internal/adapters/http/handlers/category_handler.go` | `internal/adapters/http/handlers/category/category_handler.go` |
| `internal/adapters/http/handlers/product_handler.go` | `internal/adapters/http/handlers/product/product_handler.go` |
| `internal/adapters/http/handlers/customer_handler.go` | `internal/adapters/http/handlers/customer/customer_handler.go` |
| `internal/adapters/http/handlers/sale_handler.go` | `internal/adapters/http/handlers/sale/sale_handler.go` |
| `internal/adapters/http/handlers/supplier_handler.go` | `internal/adapters/http/handlers/supplier/supplier_handler.go` |
| `internal/adapters/http/handlers/purchase_handler.go` | `internal/adapters/http/handlers/purchase/purchase_handler.go` |

### 2. **Servicios Antiguos** (carpeta completa)

La carpeta `services/` ya no se usa. Ahora usamos **casos de uso**.

| ❌ Carpeta Antigua | ✅ Carpeta Nueva |
|-------------------|------------------|
| `internal/application/services/` | `internal/application/usecases/` |

**Archivos dentro:**
- `auth_service.go`
- `user_service.go`
- `capital_injection_service.go`
- `category_service.go`
- `product_service.go`
- `customer_service.go`
- `sale_service.go`
- `supplier_service.go`
- `purchase_service.go`
- `file_service.go`

### 3. **Repositorios Antiguos** (carpeta postgres/)

Los repositorios ahora están organizados por entidad.

| ❌ Carpeta Antigua | ✅ Carpeta Nueva |
|-------------------|------------------|
| `internal/adapters/postgres/` | `internal/adapters/persistence/repositories/` |

**Archivos dentro:**
- `user_repository.go`
- `capital_injection_repository.go`
- `category_repository.go`
- `product_repository.go`
- `customer_repository.go`
- `sale_repository.go`
- `supplier_repository.go`
- `purchase_repository.go`

### 4. **Interfaces Antiguas** (ports/)

Las interfaces ahora están en `domain/ports/`.

| ❌ Archivo Antiguo | ✅ Archivo Nuevo |
|-------------------|------------------|
| `internal/ports/repositories.go` | `internal/domain/ports/*_repository.go` |
| `internal/ports/services.go` | `internal/domain/ports/services.go` |

### 5. **Entidades de Dominio Antiguas** (con GORM)

Las entidades ahora están limpias (sin GORM) en `domain/entities/`.

| ❌ Archivo Antiguo (con GORM) | ✅ Archivo Nuevo (sin GORM) |
|-------------------------------|----------------------------|
| `internal/domain/user.go` | `internal/domain/entities/user.go` |
| `internal/domain/capital_injection.go` | `internal/domain/entities/capital_injection.go` |
| `internal/domain/category.go` | `internal/domain/entities/category.go` |
| `internal/domain/product.go` | `internal/domain/entities/product.go` |
| `internal/domain/customer.go` | `internal/domain/entities/customer.go` |
| `internal/domain/sale.go` | `internal/domain/entities/sale.go` |
| `internal/domain/supplier.go` | `internal/domain/entities/supplier.go` |
| `internal/domain/purchase.go` | `internal/domain/entities/purchase.go` |

**Los modelos GORM ahora están en:**
- `internal/adapters/persistence/models/*_model.go`

---

## 🚀 Cómo Limpiar

### Opción 1: Script Automático (Recomendado)

```bash
# El script crea un backup antes de eliminar
./scripts/cleanup-old-files.sh
```

**Características:**
- ✅ Crea backup automático
- ✅ Muestra lista de archivos antes de eliminar
- ✅ Requiere confirmación
- ✅ Permite restaurar si algo sale mal

### Opción 2: Manual

Si prefieres hacerlo manualmente:

```bash
# 1. Crear backup
mkdir -p backup_manual
cp -r internal backup_manual/

# 2. Eliminar handlers antiguos
rm internal/adapters/http/handlers/*_handler.go

# 3. Eliminar servicios
rm -rf internal/application/services

# 4. Eliminar repositorios antiguos
rm -rf internal/adapters/postgres

# 5. Eliminar ports antiguos
rm -rf internal/ports

# 6. Eliminar domain antiguos
rm internal/domain/*.go
```

---

## ✅ Verificar que Todo Funciona

Después de limpiar, verifica que el proyecto compile:

```bash
# 1. Limpiar caché de Go
go clean -cache

# 2. Actualizar dependencias
go mod tidy

# 3. Compilar
go build -o /dev/null ./cmd/api/main.go

# 4. Ejecutar tests (si existen)
go test ./...
```

Si todo compila correctamente, ¡la limpieza fue exitosa! 🎉

---

## 🔄 Restaurar desde Backup

Si algo sale mal:

```bash
# Restaurar desde el backup automático
cp -r backup_old_files_*/internal/* internal/

# O desde backup manual
cp -r backup_manual/internal/* internal/
```

---

## 📊 Estructura Final

Después de la limpieza, tu estructura será:

```
internal/
├── domain/
│   ├── entities/          ✅ Entidades puras (sin GORM)
│   │   ├── user.go
│   │   ├── product.go
│   │   └── ...
│   └── ports/             ✅ Interfaces
│       ├── user_repository.go
│       ├── product_repository.go
│       └── ...
├── application/
│   └── usecases/          ✅ Casos de uso por entidad
│       ├── user/
│       ├── product/
│       └── ...
└── adapters/
    ├── persistence/
    │   ├── models/        ✅ Modelos GORM
    │   │   ├── user_model.go
    │   │   └── ...
    │   └── repositories/  ✅ Implementaciones
    │       ├── user/
    │       ├── product/
    │       └── ...
    └── http/
        ├── handlers/      ✅ Handlers por entidad
        │   ├── user/
        │   ├── product/
        │   └── ...
        ├── middleware/
        └── routes/
```

---

## 📝 Checklist

Antes de limpiar:
- [ ] El proyecto compila sin errores
- [ ] Todos los tests pasan (si existen)
- [ ] Has hecho commit de los cambios actuales
- [ ] Tienes un backup

Después de limpiar:
- [ ] El proyecto sigue compilando
- [ ] Los tests siguen pasando
- [ ] La aplicación funciona correctamente
- [ ] No hay imports rotos

---

## 💡 Consejos

1. **Haz commit antes de limpiar:**
   ```bash
   git add .
   git commit -m "refactor: complete hexagonal architecture"
   ```

2. **Verifica imports:**
   ```bash
   # Buscar imports antiguos
   grep -r "internal/ports" internal/
   grep -r "internal/application/services" internal/
   grep -r "internal/adapters/postgres" internal/
   ```

3. **Limpia imports no usados:**
   ```bash
   goimports -w .
   ```

---

¿Listo para limpiar? Ejecuta:
```bash
./scripts/cleanup-old-files.sh
```
