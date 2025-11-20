# 🚀 Quick Start - Fashion Blue

## ✅ Servicios Levantados

Los siguientes servicios están corriendo:

| Servicio | Estado | URL | Credenciales |
|----------|--------|-----|--------------|
| **PostgreSQL** | ✅ Running | `localhost:5432` | User: `fashionblue`<br>Pass: `fashionblue123`<br>DB: `fashionblue_db` |
| **pgAdmin** | ✅ Running | http://localhost:5050 | Email: `admin@fashionblue.com`<br>Pass: `admin123` |

---

## 🎯 Siguiente Paso: Ejecutar la Aplicación

Tienes 3 opciones:

### Opción 1: Ejecutar con Go Run (Simple)
```bash
go run cmd/api/main.go
```

### Opción 2: Ejecutar con Air (Hot Reload)
```bash
# Instalar air si no lo tienes
go install github.com/cosmtrek/air@latest

# Ejecutar
air
```

### Opción 3: Debugging con VS Code (Recomendado)

1. **Copiar configuración de debugging:**
   ```bash
   mkdir -p .vscode
   cp launch.json.example .vscode/launch.json
   ```

2. **En VS Code:**
   - Presiona `F5`
   - O ve a "Run and Debug" (Ctrl+Shift+D)
   - Selecciona "Debug Fashion Blue API"
   - ¡Listo! Puedes poner breakpoints

---

## 🧪 Probar que Funciona

Una vez que la app esté corriendo, prueba:

### 1. Health Check
```bash
curl http://localhost:8080/health
```

Deberías ver:
```json
{
  "status": "healthy",
  "time": "2024-11-20T12:45:00Z"
}
```

### 2. Registrar un Usuario
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "admin123",
    "first_name": "Admin",
    "last_name": "User",
    "role": "SUPER_ADMIN"
  }'
```

### 3. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "admin123"
  }'
```

Guarda el `token` que te devuelve para usarlo en las siguientes peticiones.

---

## 📊 Acceder a pgAdmin

1. Abre http://localhost:5050
2. Login:
   - Email: `admin@fashionblue.com`
   - Password: `admin123`
3. Agregar servidor:
   - Click derecho en "Servers" → "Register" → "Server"
   - **General Tab:**
     - Name: `Fashion Blue DB`
   - **Connection Tab:**
     - Host: `postgres`
     - Port: `5432`
     - Username: `fashionblue`
     - Password: `fashionblue123`
     - Save password: ✅
   - Click "Save"

---

## 🛑 Detener los Servicios

```bash
# Detener solo los contenedores
docker-compose down

# Detener y eliminar volúmenes (⚠️ borra los datos)
docker-compose down -v
```

---

## 📚 Más Información

- Ver `DEBUG_GUIDE.md` para guía completa de debugging
- Ver `API_EXAMPLES.md` para ejemplos de todos los endpoints
- Ver `README.md` para documentación general

---

## ⚡ Comandos Útiles

```bash
# Ver logs de PostgreSQL
docker-compose logs -f postgres

# Ver estado de los contenedores
docker ps

# Reiniciar PostgreSQL
docker-compose restart postgres

# Conectar a PostgreSQL desde terminal
psql -h localhost -U fashionblue -d fashionblue_db

# Ver tablas en la base de datos
docker exec -it fashionblue-postgres psql -U fashionblue -d fashionblue_db -c "\dt"
```

---

¡Todo listo para empezar a desarrollar! 🎉
