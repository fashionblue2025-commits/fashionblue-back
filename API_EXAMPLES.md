# Fashion Blue - Ejemplos de Uso de la API (Actualizado)

Este documento contiene ejemplos prácticos de cómo usar la API de Fashion Blue con el sistema simplificado de transacciones manuales.

## 🔐 Autenticación

### 1. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@fashionblue.com",
    "password": "admin123"
  }'
```

**Respuesta:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "admin@fashionblue.com",
      "firstName": "Super",
      "lastName": "Admin",
      "role": "SUPER_ADMIN",
      "isActive": true
    }
  }
}
```

### 2. Registrar nuevo usuario (vendedor)

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "vendedor@fashionblue.com",
    "password": "vendedor123",
    "first_name": "Juan",
    "last_name": "Pérez",
    "role": "SELLER"
  }'
```

## 💰 Inyecciones de Capital

### Crear inyección de capital (Solo Super Admin)

```bash
curl -X POST http://localhost:8080/api/v1/capital-injections \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d '{
    "amount": 5000000,
    "description": "Inversión inicial para compra de materiales",
    "source": "Inversión personal",
    "date": "2024-01-15T00:00:00Z"
  }'
```

### Listar inyecciones de capital

```bash
curl -X GET http://localhost:8080/api/v1/capital-injections \
  -H "Authorization: Bearer TU_TOKEN"
```

### Obtener total de capital

```bash
curl -X GET http://localhost:8080/api/v1/capital-injections/total \
  -H "Authorization: Bearer TU_TOKEN"
```

## 🏭 Proveedores

### Crear proveedor (Solo Super Admin)

```bash
curl -X POST http://localhost:8080/api/v1/suppliers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d '{
    "name": "Textiles del Norte S.A.S",
    "contact_name": "Carlos Ramírez",
    "email": "ventas@textilesnorte.com",
    "phone": "3201234567",
    "address": "Zona Industrial Calle 80, Barranquilla",
    "notes": "Proveedor principal de telas",
    "is_active": true
  }'
```

### Listar proveedores

```bash
curl -X GET http://localhost:8080/api/v1/suppliers \
  -H "Authorization: Bearer TU_TOKEN"
```

### Actualizar proveedor

```bash
curl -X PUT http://localhost:8080/api/v1/suppliers/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d '{
    "name": "Textiles del Norte S.A.S",
    "contact_name": "Carlos Ramírez Gómez",
    "phone": "3201234567",
    "address": "Nueva dirección",
    "is_active": true
  }'
```

## 👥 Clientes y Transacciones Manuales

### Crear cliente

```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d '{
    "name": "María González",
    "phone": "3001234567",
    "address": "Calle 123 #45-67, Bogotá",
    "risk_level": "LOW",
    "payment_frequency": "BIWEEKLY",
    "payment_days": "2,17",
    "notes": "Cliente frecuente"
  }'
```

### Registrar DEUDA (cuando el cliente compra)

```bash
curl -X POST http://localhost:8080/api/v1/customers/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d '{
    "customer_id": 1,
    "type": "DEUDA",
    "amount": 150000,
    "concept": "Venta de chaqueta de cuero",
    "date": "2024-11-20T00:00:00Z"
  }'
```

### Registrar ABONO (cuando el cliente paga)

```bash
curl -X POST http://localhost:8080/api/v1/customers/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d '{
    "customer_id": 1,
    "type": "ABONO",
    "amount": 50000,
    "payment_method_id": 1,
    "concept": "Abono quincenal",
    "date": "2024-11-20T00:00:00Z"
  }'
```

### Ver balance de un cliente

```bash
curl -X GET http://localhost:8080/api/v1/customers/1/balance \
  -H "Authorization: Bearer TU_TOKEN"
```

**Respuesta:**
```json
{
  "success": true,
  "message": "Balance retrieved successfully",
  "data": {
    "customerId": 1,
    "balance": 100000
  }
}
```

### Ver historial de transacciones

```bash
curl -X GET http://localhost:8080/api/v1/customers/1/history \
  -H "Authorization: Bearer TU_TOKEN"
```

### Listar clientes próximos a pagar

```bash
curl -X GET http://localhost:8080/api/v1/customers/upcoming-payments \
  -H "Authorization: Bearer TU_TOKEN"
```

## 📦 Productos

### Crear producto

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d '{
    "name": "Chaqueta de Cuero Premium",
    "description": "Chaqueta de cuero genuino",
    "sku": "CHQ-CUERO-001",
    "category_id": 1,
    "material_cost": 80000,
    "labor_cost": 50000,
    "unit_price": 250000,
    "stock": 50,
    "min_stock": 10
  }'
```

### Listar productos

```bash
curl -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer TU_TOKEN"
```

### Productos con stock bajo

```bash
curl -X GET http://localhost:8080/api/v1/products/low-stock \
  -H "Authorization: Bearer TU_TOKEN"
```

## 💡 Caso de Uso Completo: Venta a Crédito

```bash
# 1. Login
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@fashionblue.com","password":"admin123"}' \
  | jq -r '.data.token')

# 2. Crear cliente
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "María González",
    "phone": "3001234567",
    "address": "Calle 123, Bogotá",
    "risk_level": "LOW",
    "payment_frequency": "BIWEEKLY",
    "payment_days": "2,17"
  }'

# 3. Registrar DEUDA por la venta
curl -X POST http://localhost:8080/api/v1/customers/transactions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "type": "DEUDA",
    "amount": 150000,
    "concept": "Venta chaqueta de cuero"
  }'

# 4. Ver balance
curl -X GET http://localhost:8080/api/v1/customers/1/balance \
  -H "Authorization: Bearer $TOKEN"
# Balance: 150000

# 5. Cliente hace primer abono
curl -X POST http://localhost:8080/api/v1/customers/transactions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "type": "ABONO",
    "amount": 50000,
    "payment_method_id": 1,
    "concept": "Primer abono"
  }'

# 6. Ver balance actualizado
curl -X GET http://localhost:8080/api/v1/customers/1/balance \
  -H "Authorization: Bearer $TOKEN"
# Balance: 100000 (150000 - 50000)
```

## 📝 Notas Importantes

### Sistema de Transacciones Manuales
- **DEUDA**: Se registra cuando el cliente compra algo (aumenta el balance)
- **ABONO**: Se registra cuando el cliente paga (disminuye el balance)
- El balance se calcula automáticamente sumando todas las transacciones

### Roles y Permisos
- **SUPER_ADMIN**: Acceso completo, puede crear proveedores e inyecciones de capital
- **SELLER**: Puede registrar transacciones de clientes, ver productos

### Endpoints Principales
- `/api/v1/auth/*` - Autenticación (público)
- `/api/v1/capital-injections/*` - Inyecciones de capital (SuperAdmin)
- `/api/v1/suppliers/*` - Proveedores (SuperAdmin para crear/editar)
- `/api/v1/customers/*` - Clientes y transacciones
- `/api/v1/products/*` - Productos
- `/api/v1/categories/*` - Categorías
- `/api/v1/payment-methods/*` - Métodos de pago

### Seguridad
- Todos los endpoints requieren autenticación excepto `/auth/login` y `/auth/register`
- Token JWT debe enviarse en header: `Authorization: Bearer TOKEN`
- Tokens expiran en 24 horas

## 🚨 Códigos de Estado HTTP

- `200 OK` - Solicitud exitosa
- `201 Created` - Recurso creado
- `400 Bad Request` - Error en datos
- `401 Unauthorized` - No autenticado
- `403 Forbidden` - Sin permisos
- `404 Not Found` - No encontrado
- `500 Internal Server Error` - Error del servidor
