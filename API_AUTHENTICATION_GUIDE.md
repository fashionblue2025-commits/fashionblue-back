# 🔐 Guía de Autenticación - Fashion Blue API

## 📋 Resumen

La API usa **JWT (JSON Web Tokens)** para autenticación. El usuario autenticado se obtiene automáticamente del token en las rutas protegidas.

---

## 🔑 Obtener Token

### 1. Registrar Usuario

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@fashionblue.com",
    "password": "Admin123!",
    "first_name": "Admin",
    "last_name": "User",
    "role": "SUPER_ADMIN"
  }'
```

**Respuesta:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "admin@fashionblue.com",
      "first_name": "Admin",
      "last_name": "User",
      "role": "SUPER_ADMIN"
    }
  }
}
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@fashionblue.com",
    "password": "Admin123!"
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
      "first_name": "Admin",
      "last_name": "User",
      "role": "SUPER_ADMIN"
    }
  }
}
```

---

## 🛡️ Usar el Token

### Guardar el Token

```bash
# Guardar en variable de entorno
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# O en archivo
echo "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." > .token
```

### Hacer Peticiones Autenticadas

```bash
# Usando variable de entorno
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"

# Usando archivo
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $(cat .token)"

# Directamente
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 🎯 Campos Automáticos del Usuario Autenticado

Los siguientes campos se obtienen **automáticamente** del token JWT y **NO deben enviarse** en el body:

### ✅ Inyección de Capital

**❌ ANTES (Incorrecto):**
```json
{
  "amount": 5000000,
  "type": "MATERIALS",
  "description": "Compra de telas",
  "date": "2024-01-15T00:00:00Z",
  "created_by": 1  // ❌ NO enviar esto
}
```

**✅ AHORA (Correcto):**
```json
{
  "amount": 5000000,
  "type": "MATERIALS",
  "description": "Compra de telas",
  "date": "2024-01-15T00:00:00Z"
  // created_by se obtiene automáticamente del token
}
```

**Ejemplo completo:**
```bash
curl -X POST http://localhost:8080/api/v1/capital-injections \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "amount": 5000000,
    "type": "MATERIALS",
    "description": "Compra de telas para producción",
    "date": "2024-01-15T00:00:00Z"
  }'
```

### ✅ Compras (Purchase)

**❌ ANTES (Incorrecto):**
```json
{
  "supplier_id": 1,
  "total": 1500000,
  "notes": "Compra de materiales",
  "purchase_date": "2024-01-15T00:00:00Z",
  "created_by": 1,  // ❌ NO enviar esto
  "items": [...]
}
```

**✅ AHORA (Correcto):**
```json
{
  "supplier_id": 1,
  "total": 1500000,
  "notes": "Compra de materiales",
  "purchase_date": "2024-01-15T00:00:00Z",
  "items": [
    {
      "product_id": 1,
      "quantity": 100,
      "unit_price": 15000
    }
  ]
  // created_by se obtiene automáticamente del token
}
```

**Ejemplo completo:**
```bash
curl -X POST http://localhost:8080/api/v1/purchases \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "supplier_id": 1,
    "total": 1500000,
    "notes": "Compra mensual de telas",
    "purchase_date": "2024-01-15T00:00:00Z",
    "items": [
      {
        "product_id": 1,
        "quantity": 100,
        "unit_price": 15000
      },
      {
        "product_id": 2,
        "quantity": 50,
        "unit_price": 30000
      }
    ]
  }'
```

---

## 🔍 Información del Token JWT

El token contiene la siguiente información:

```json
{
  "user_id": 1,
  "email": "admin@fashionblue.com",
  "role": "SUPER_ADMIN",
  "exp": 1700000000,  // Expiración
  "iat": 1699913600   // Fecha de emisión
}
```

Esta información se usa para:
- ✅ Identificar al usuario (`user_id`)
- ✅ Verificar permisos (`role`)
- ✅ Validar expiración (`exp`)
- ✅ Asignar automáticamente `created_by`

---

## 🚫 Errores Comunes

### 1. Token No Enviado

```bash
curl -X POST http://localhost:8080/api/v1/capital-injections \
  -H "Content-Type: application/json" \
  -d '{...}'
```

**Error:**
```json
{
  "success": false,
  "message": "Missing authorization header"
}
```

**Solución:** Agregar header `Authorization: Bearer TOKEN`

### 2. Token Inválido

```bash
curl -X POST http://localhost:8080/api/v1/capital-injections \
  -H "Authorization: Bearer token-invalido" \
  -d '{...}'
```

**Error:**
```json
{
  "success": false,
  "message": "Invalid or expired token"
}
```

**Solución:** Hacer login nuevamente para obtener un token válido

### 3. Token Expirado

**Error:**
```json
{
  "success": false,
  "message": "Invalid or expired token"
}
```

**Solución:** Hacer login nuevamente (los tokens expiran después de 24h por defecto)

### 4. Usuario Inactivo

**Error:**
```json
{
  "success": false,
  "message": "User is inactive"
}
```

**Solución:** Contactar al administrador para activar la cuenta

---

## 🎭 Roles y Permisos

### Roles Disponibles

- **SUPER_ADMIN**: Acceso total
- **SELLER**: Puede crear ventas, ver productos, clientes
- **VIEWER**: Solo lectura

### Endpoints por Rol

| Endpoint | SUPER_ADMIN | SELLER | VIEWER |
|----------|-------------|--------|--------|
| POST /capital-injections | ✅ | ❌ | ❌ |
| POST /purchases | ✅ | ❌ | ❌ |
| POST /sales | ✅ | ✅ | ❌ |
| GET /sales | ✅ | ✅ | ✅ |
| POST /users | ✅ | ❌ | ❌ |
| GET /users | ✅ | ❌ | ❌ |

---

## 💡 Mejores Prácticas

### 1. Almacenar el Token de Forma Segura

```javascript
// ✅ BIEN - En memoria o localStorage con cuidado
localStorage.setItem('token', token);

// ❌ MAL - En cookies sin httpOnly
document.cookie = `token=${token}`;
```

### 2. Renovar Token Antes de Expirar

```javascript
// Verificar expiración
const payload = JSON.parse(atob(token.split('.')[1]));
const expiresIn = payload.exp * 1000 - Date.now();

if (expiresIn < 3600000) { // Menos de 1 hora
  // Renovar token
  refreshToken();
}
```

### 3. Limpiar Token al Cerrar Sesión

```javascript
// Eliminar token
localStorage.removeItem('token');
sessionStorage.clear();
```

### 4. Manejar Errores de Autenticación

```javascript
// Interceptor de respuestas
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // Redirigir a login
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

---

## 🔧 Configuración del Token

En el archivo `.env`:

```env
# Clave secreta (mínimo 32 caracteres)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Tiempo de expiración
JWT_EXPIRATION=24h  # 24 horas
# JWT_EXPIRATION=1h   # 1 hora
# JWT_EXPIRATION=7d   # 7 días
```

---

## 📚 Recursos Adicionales

- [JWT.io](https://jwt.io/) - Decodificar y verificar tokens
- [Postman](https://www.postman.com/) - Cliente API con soporte JWT
- [Thunder Client](https://www.thunderclient.com/) - Extension de VS Code

---

**Última actualización:** 2024-11-20
