# 👤 Guía de Refactorización - Modelo de Clientes

## 📋 Resumen de Cambios

El modelo de `Customer` ha sido completamente refactorizado para ser más flexible y realista:

### ❌ Modelo Anterior (Incorrecto)

```go
type Customer struct {
    FirstName    string
    LastName     string
    Email        string
    DocumentType string
    DocumentNum  string
    Type         CustomerType // RETAIL o WHOLESALE ❌
    // ...
}
```

**Problemas:**
- ❌ El tipo (retail/wholesale) estaba en el cliente
- ❌ Un cliente no puede comprar ambos tipos
- ❌ No había información de tallas
- ❌ No había nivel de riesgo

### ✅ Modelo Nuevo (Correcto)

```go
type Customer struct {
    ID           uint
    Name         string      // Nombre completo
    Phone        string      // Teléfono
    Address      string      // Dirección
    RiskLevel    RiskLevel   // LOW, MEDIUM, HIGH
    ShirtSizeID  *uint       // ID de talla de camiseta (opcional)
    PantsSizeID  *uint       // ID de talla de pantalón (opcional)
    ShoesSizeID  *uint       // ID de talla de tenis (opcional)
    Birthday     *time.Time  // Fecha de cumpleaños (opcional)
    Notes        string      // Notas adicionales
    IsActive     bool
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

**Ventajas:**
- ✅ El tipo de venta (retail/wholesale) ahora está en la **venta**, no en el cliente
- ✅ Un cliente puede comprar al detal y al por mayor
- ✅ Información de tallas para personalización
- ✅ Nivel de riesgo para gestión de crédito
- ✅ Cumpleaños para marketing
- ✅ Más simple y flexible

---

## 🎯 Nueva Entidad: Size (Talla)

```go
type SizeType string

const (
    SizeTypeShirt SizeType = "SHIRT"  // Camiseta
    SizeTypePants SizeType = "PANTS"  // Pantalón
    SizeTypeShoes SizeType = "SHOES"  // Tenis/Zapatos
)

type Size struct {
    ID        uint
    Type      SizeType  // SHIRT, PANTS, SHOES
    Value     string    // "S", "M", "L", "28", "30", "7", "8.5", etc.
    Order     int       // Para ordenar (S=1, M=2, L=3, etc.)
    IsActive  bool
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Tallas Predefinidas

#### Camisetas (SHIRT)
- XS, S, M, L, XL, XXL

#### Pantalones (PANTS)
- 24, 26, 28, 30, 32, 34, 36, 38, 40, 42 (pulgadas)

#### Zapatos/Tenis (SHOES)
- 5, 5.5, 6, 6.5, 7, 7.5, 8, 8.5, 9, 9.5, 10, 10.5, 11, 11.5, 12, 13, 14 (US)

---

## 🔄 Migración de Datos

### 1. Ejecutar Migraciones

```bash
# La migración se ejecuta automáticamente al iniciar la app
go run cmd/api/main.go
```

Esto creará:
- ✅ Tabla `sizes` (nueva)
- ✅ Tabla `customers` (actualizada con nuevas columnas)

### 2. Poblar Tallas

```bash
# Ejecutar script de seed para tallas
go run scripts/seed_sizes.go
```

Esto insertará:
- 6 tallas de camisetas
- 10 tallas de pantalones
- 17 tallas de zapatos
- **Total: 33 tallas**

### 3. Migrar Clientes Existentes (Si hay datos)

Si ya tienes clientes en la BD, necesitas migrarlos:

```sql
-- Actualizar clientes existentes con valores por defecto
UPDATE customers 
SET 
    name = CONCAT(first_name, ' ', last_name),
    risk_level = 'LOW',
    is_active = true
WHERE name IS NULL OR name = '';

-- Opcional: Eliminar columnas antiguas (después de verificar)
ALTER TABLE customers 
DROP COLUMN IF EXISTS first_name,
DROP COLUMN IF EXISTS last_name,
DROP COLUMN IF EXISTS email,
DROP COLUMN IF EXISTS city,
DROP COLUMN IF EXISTS document_type,
DROP COLUMN IF EXISTS document_num,
DROP COLUMN IF EXISTS type;
```

---

## 📝 Nuevos Endpoints de API

### Listar Tallas

```bash
# Todas las tallas
curl -X GET http://localhost:8080/api/v1/sizes \
  -H "Authorization: Bearer $TOKEN"

# Tallas de camisetas
curl -X GET "http://localhost:8080/api/v1/sizes?type=SHIRT" \
  -H "Authorization: Bearer $TOKEN"

# Tallas de pantalones
curl -X GET "http://localhost:8080/api/v1/sizes?type=PANTS" \
  -H "Authorization: Bearer $TOKEN"

# Tallas de zapatos
curl -X GET "http://localhost:8080/api/v1/sizes?type=SHOES" \
  -H "Authorization: Bearer $TOKEN"
```

### Crear Cliente (Nuevo Formato)

```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Juan Pérez",
    "phone": "3001234567",
    "address": "Calle 123 #45-67, Bogotá",
    "risk_level": "LOW",
    "shirt_size_id": 3,
    "pants_size_id": 5,
    "shoes_size_id": 10,
    "birthday": "1990-05-15T00:00:00Z",
    "notes": "Cliente frecuente, prefiere colores oscuros"
  }'
```

**Campos Opcionales:**
- `shirt_size_id`: ID de la talla de camiseta
- `pants_size_id`: ID de la talla de pantalón
- `shoes_size_id`: ID de la talla de zapatos
- `birthday`: Fecha de cumpleaños
- `notes`: Notas adicionales

**Respuesta:**
```json
{
  "success": true,
  "message": "Customer created successfully",
  "data": {
    "id": 1,
    "name": "Juan Pérez",
    "phone": "3001234567",
    "address": "Calle 123 #45-67, Bogotá",
    "risk_level": "LOW",
    "shirt_size_id": 3,
    "pants_size_id": 5,
    "shoes_size_id": 10,
    "birthday": "1990-05-15T00:00:00Z",
    "notes": "Cliente frecuente, prefiere colores oscuros",
    "is_active": true,
    "shirt_size": {
      "id": 3,
      "type": "SHIRT",
      "value": "M",
      "order": 3
    },
    "pants_size": {
      "id": 5,
      "type": "PANTS",
      "value": "32",
      "order": 5
    },
    "shoes_size": {
      "id": 10,
      "type": "SHOES",
      "value": "9.5",
      "order": 10
    }
  }
}
```

---

## 🎨 Casos de Uso

### 1. Cliente con Tallas Conocidas

```json
{
  "name": "María García",
  "phone": "3009876543",
  "address": "Carrera 7 #100-50",
  "risk_level": "LOW",
  "shirt_size_id": 2,
  "pants_size_id": 3,
  "shoes_size_id": 7,
  "birthday": "1995-08-20T00:00:00Z"
}
```

**Beneficio:** Puedes recomendar productos de su talla automáticamente.

### 2. Cliente sin Tallas (Aún no las conoces)

```json
{
  "name": "Carlos Rodríguez",
  "phone": "3001112233",
  "address": "Calle 50 #20-30",
  "risk_level": "MEDIUM"
}
```

**Beneficio:** Puedes agregar las tallas después cuando las conozcas.

### 3. Cliente de Alto Riesgo

```json
{
  "name": "Pedro López",
  "phone": "3004445566",
  "address": "Avenida 1 #2-3",
  "risk_level": "HIGH",
  "notes": "Historial de pagos atrasados"
}
```

**Beneficio:** Puedes aplicar políticas especiales (pago anticipado, límite de crédito, etc.).

### 4. Marketing de Cumpleaños

```go
// Obtener clientes que cumplen años hoy
customers := customerRepo.GetBirthdayCustomers(time.Now())

for _, customer := range customers {
    if customer.IsBirthday() {
        // Enviar mensaje de cumpleaños
        // Ofrecer descuento especial
        sendBirthdayPromotion(customer)
    }
}
```

---

## 🔄 Actualización de Ventas

El tipo de venta (retail/wholesale) ahora está en la **venta**, no en el cliente:

```go
type Sale struct {
    ID         uint
    CustomerID uint
    Type       SaleType      // RETAIL o WHOLESALE ✅
    Total      float64
    // ...
}
```

**Ejemplo:**
```bash
# Venta al detal
curl -X POST http://localhost:8080/api/v1/sales \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "customer_id": 1,
    "type": "RETAIL",
    "items": [...]
  }'

# Venta al por mayor (mismo cliente)
curl -X POST http://localhost:8080/api/v1/sales \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "customer_id": 1,
    "type": "WHOLESALE",
    "items": [...]
  }'
```

---

## 📊 Reportes y Análisis

### Clientes por Nivel de Riesgo

```sql
SELECT risk_level, COUNT(*) as total
FROM customers
WHERE is_active = true
GROUP BY risk_level;
```

### Tallas Más Comunes

```sql
SELECT s.type, s.value, COUNT(c.id) as customers
FROM sizes s
LEFT JOIN customers c ON (
    s.id = c.shirt_size_id OR 
    s.id = c.pants_size_id OR 
    s.id = c.shoes_size_id
)
GROUP BY s.type, s.value
ORDER BY customers DESC;
```

### Cumpleaños del Mes

```sql
SELECT name, phone, birthday
FROM customers
WHERE EXTRACT(MONTH FROM birthday) = EXTRACT(MONTH FROM CURRENT_DATE)
  AND is_active = true
ORDER BY EXTRACT(DAY FROM birthday);
```

---

## ✅ Checklist de Implementación

- [x] Crear entidad `Size`
- [x] Actualizar entidad `Customer`
- [x] Crear modelo de persistencia `SizeModel`
- [x] Actualizar modelo de persistencia `CustomerModel`
- [x] Agregar migración de tabla `sizes`
- [x] Crear script de seed para tallas
- [ ] Crear repositorio de `Size`
- [ ] Crear casos de uso de `Size`
- [ ] Crear handler de `Size`
- [ ] Actualizar handler de `Customer`
- [ ] Actualizar entidad `Sale` con tipo de venta
- [ ] Actualizar documentación de API
- [ ] Migrar datos existentes (si hay)
- [ ] Actualizar tests

---

## 🚀 Próximos Pasos

1. **Ejecutar migraciones:**
   ```bash
   go run cmd/api/main.go
   ```

2. **Poblar tallas:**
   ```bash
   go run scripts/seed_sizes.go
   ```

3. **Probar API:**
   ```bash
   # Listar tallas
   curl http://localhost:8080/api/v1/sizes

   # Crear cliente con tallas
   curl -X POST http://localhost:8080/api/v1/customers \
     -H "Authorization: Bearer $TOKEN" \
     -d '{...}'
   ```

---

**Última actualización:** 2024-11-20
