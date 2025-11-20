# 🐛 Guía de Debugging - Fashion Blue

## 🚀 Inicio Rápido

### Opción 1: Desarrollo Local (Recomendado para debugging)

```bash
# 1. Dar permisos al script
chmod +x dev.sh

# 2. Levantar solo la base de datos
./dev.sh
```

Esto levantará:
- ✅ PostgreSQL en `localhost:5432`
- ✅ pgAdmin en `http://localhost:5050`

Luego puedes ejecutar la app de 3 formas:

#### A. Con Go Run
```bash
go run cmd/api/main.go
```

#### B. Con Air (hot reload)
```bash
# Instalar air si no lo tienes
go install github.com/cosmtrek/air@latest

# Ejecutar con hot reload
air
```

#### C. Con el Debugger de VS Code
1. Copia `launch.json.example` a `.vscode/launch.json`
2. Presiona `F5` o ve a "Run and Debug"
3. Selecciona "Debug Fashion Blue API"
4. ¡Listo! Puedes poner breakpoints

---

### Opción 2: Todo en Docker

```bash
# Levantar todos los servicios
docker-compose up -d

# Ver logs
docker-compose logs -f api

# Detener
docker-compose down
```

---

## 📊 Servicios Disponibles

| Servicio | URL | Credenciales |
|----------|-----|--------------|
| **API** | http://localhost:8080 | - |
| **Health Check** | http://localhost:8080/health | - |
| **PostgreSQL** | localhost:5432 | User: `fashionblue`<br>Pass: `fashionblue123`<br>DB: `fashionblue_db` |
| **pgAdmin** | http://localhost:5050 | Email: `admin@fashionblue.com`<br>Pass: `admin123` |

---

## 🔧 Configurar pgAdmin

1. Abre http://localhost:5050
2. Login con `admin@fashionblue.com` / `admin123`
3. Click derecho en "Servers" → "Register" → "Server"
4. En la pestaña "General":
   - Name: `Fashion Blue DB`
5. En la pestaña "Connection":
   - Host: `postgres` (si usas Docker) o `localhost` (si usas dev.sh)
   - Port: `5432`
   - Username: `fashionblue`
   - Password: `fashionblue123`
   - Save password: ✅
6. Click "Save"

---

## 🐛 Debugging con VS Code

### 1. Configuración Inicial

```bash
# Copiar configuración de debugging
cp launch.json.example .vscode/launch.json
```

### 2. Usar el Debugger

1. **Levantar la base de datos:**
   ```bash
   ./dev.sh
   ```

2. **Abrir VS Code** en el proyecto

3. **Poner breakpoints** en el código (click en el margen izquierdo)

4. **Iniciar debugging:**
   - Presiona `F5`
   - O ve a "Run and Debug" (Ctrl+Shift+D)
   - Selecciona "Debug Fashion Blue API"
   - Click en el botón verde "Start Debugging"

5. **Controles del debugger:**
   - `F5` - Continue
   - `F10` - Step Over
   - `F11` - Step Into
   - `Shift+F11` - Step Out
   - `Shift+F5` - Stop

### 3. Variables de Entorno

El debugger usa las variables definidas en `launch.json`. Si necesitas cambiarlas:
- Edita `.vscode/launch.json`
- O crea un archivo `.env` (tiene prioridad)

---

## 🧪 Probar la API

### Health Check
```bash
curl http://localhost:8080/health
```

### Registrar Usuario
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

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "admin123"
  }'
```

---

## 📝 Comandos Útiles

### Docker
```bash
# Ver logs de la API
docker-compose logs -f api

# Ver logs de PostgreSQL
docker-compose logs -f postgres

# Reiniciar un servicio
docker-compose restart api

# Entrar al contenedor de PostgreSQL
docker exec -it fashionblue-postgres psql -U fashionblue -d fashionblue_db

# Limpiar todo (⚠️ borra los datos)
docker-compose down -v
```

### Base de Datos
```bash
# Conectar a PostgreSQL desde la terminal
psql -h localhost -U fashionblue -d fashionblue_db

# Ver tablas
\dt

# Ver estructura de una tabla
\d users

# Salir
\q
```

### Go
```bash
# Compilar
go build -o bin/api cmd/api/main.go

# Ejecutar compilado
./bin/api

# Ver dependencias
go mod tidy

# Actualizar dependencias
go get -u ./...

# Ejecutar tests
go test ./...

# Ver cobertura
go test -cover ./...
```

---

## 🔍 Debugging Tips

### 1. Ver Queries SQL
En `pkg/database/postgres.go`, el logger de GORM está configurado en modo `logger.Info` en desarrollo. Verás todas las queries en la consola.

### 2. Logs Estructurados
Los logs están en formato JSON. Para verlos más legibles:
```bash
# Con jq
go run cmd/api/main.go | jq

# O cambia LOG_FORMAT=text en .env
```

### 3. Breakpoints Condicionales
En VS Code, click derecho en un breakpoint → "Edit Breakpoint" → "Expression"
```go
// Ejemplo: solo parar si el ID es 5
user.ID == 5
```

### 4. Watch Variables
En el panel de debugging, agrega variables al "Watch" para monitorearlas.

### 5. Debug Console
Puedes evaluar expresiones Go en tiempo real en la "Debug Console".

---

## ⚠️ Troubleshooting

### Puerto 5432 ya está en uso
```bash
# Ver qué está usando el puerto
lsof -i :5432

# Detener PostgreSQL local si está corriendo
brew services stop postgresql
# o
sudo systemctl stop postgresql
```

### Puerto 8080 ya está en uso
```bash
# Ver qué está usando el puerto
lsof -i :8080

# Matar el proceso
kill -9 <PID>
```

### Error de conexión a la base de datos
```bash
# Verificar que PostgreSQL está corriendo
docker ps | grep postgres

# Ver logs de PostgreSQL
docker-compose logs postgres

# Reiniciar PostgreSQL
docker-compose restart postgres
```

### Cambios no se reflejan
```bash
# Si usas Docker, reconstruir la imagen
docker-compose up -d --build api

# Si usas go run, asegúrate de que no haya procesos antiguos
pkill -f "go run"
```

---

## 📚 Recursos Adicionales

- [Documentación de Echo](https://echo.labstack.com/)
- [Documentación de GORM](https://gorm.io/docs/)
- [VS Code Go Debugging](https://github.com/golang/vscode-go/wiki/debugging)
- [Docker Compose](https://docs.docker.com/compose/)

---

## 🎯 Workflow Recomendado

1. **Levantar servicios:** `./dev.sh`
2. **Abrir VS Code**
3. **Poner breakpoints** en el código
4. **Presionar F5** para iniciar debugging
5. **Hacer requests** con curl/Postman/Thunder Client
6. **Inspeccionar variables** en el debugger
7. **Iterar** hasta resolver el issue

¡Happy Debugging! 🐛✨
