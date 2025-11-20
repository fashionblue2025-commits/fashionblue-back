# Fashion Blue - Guía de Configuración

## 📋 Requisitos Previos

- Go 1.21 o superior
- Docker y Docker Compose
- PostgreSQL (si no usas Docker)
- Make (opcional, pero recomendado)

## 🚀 Instalación Rápida

### 1. Clonar el repositorio

```bash
git clone <repository-url>
cd fashion-blue
```

### 2. Configurar variables de entorno

```bash
cp .env.example .env
```

Edita el archivo `.env` con tus configuraciones:

```env
# Application
APP_NAME=fashion-blue
APP_ENV=development
APP_PORT=8080
APP_HOST=0.0.0.0

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=fashionblue
DB_PASSWORD=fashionblue123
DB_NAME=fashionblue_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=tu-clave-secreta-muy-segura-cambiala-en-produccion
JWT_EXPIRATION=24h
```

### 3. Iniciar la base de datos con Docker

```bash
docker-compose up -d postgres
```

O si quieres iniciar todos los servicios:

```bash
docker-compose up -d
```

### 4. Instalar dependencias de Go

```bash
go mod download
```

### 5. Ejecutar migraciones y seed inicial

```bash
go run scripts/seed.go
```

Esto creará:
- Las tablas de la base de datos
- Un usuario super admin:
  - **Email**: `admin@fashionblue.com`
  - **Password**: `admin123`
- Categorías iniciales (Chaquetas, Pantalones, Camisas)

### 6. Ejecutar la aplicación

```bash
go run cmd/api/main.go
```

O usando Make:

```bash
make run
```

La API estará disponible en: `http://localhost:8080`

## 🧪 Probar la API

### Health Check

```bash
curl http://localhost:8080/health
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@fashionblue.com",
    "password": "admin123"
  }'
```

Esto te devolverá un token JWT que debes usar en las siguientes peticiones.

### Crear una categoría (requiere autenticación)

```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN_AQUI" \
  -d '{
    "name": "Zapatos",
    "description": "Calzado en general",
    "is_active": true
  }'
```

## 🐳 Docker

### Construir y ejecutar con Docker

```bash
# Construir la imagen
docker-compose build

# Iniciar todos los servicios
docker-compose up -d

# Ver logs
docker-compose logs -f api

# Detener servicios
docker-compose down
```

### Acceder a pgAdmin

Si iniciaste todos los servicios con Docker Compose, puedes acceder a pgAdmin en:

- URL: `http://localhost:5050`
- Email: `admin@fashionblue.com`
- Password: `admin123`

Para conectar a la base de datos en pgAdmin:
- Host: `postgres`
- Port: `5432`
- Database: `fashionblue_db`
- Username: `fashionblue`
- Password: `fashionblue123`

## 📝 Comandos Útiles

```bash
# Ejecutar la aplicación
make run

# Ejecutar tests
make test

# Ver cobertura de tests
make test-coverage

# Compilar la aplicación
make build

# Limpiar archivos generados
make clean

# Formatear código
make format

# Ejecutar linter
make lint

# Instalar dependencias
make deps
```

## 🔧 Desarrollo

### Estructura del Proyecto

```
fashion-blue/
├── cmd/
│   └── api/
│       └── main.go              # Punto de entrada
├── internal/
│   ├── domain/                  # Entidades de dominio
│   ├── ports/                   # Interfaces (puertos)
│   ├── application/
│   │   └── services/            # Casos de uso
│   └── adapters/
│       ├── http/                # Handlers HTTP
│       │   ├── handlers/
│       │   ├── middleware/
│       │   └── routes/
│       └── postgres/            # Repositorios PostgreSQL
├── pkg/                         # Paquetes compartidos
│   ├── config/
│   ├── database/
│   └── response/
├── scripts/                     # Scripts de utilidad
├── migrations/                  # Migraciones SQL
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

### Agregar nuevas migraciones

Si necesitas crear migraciones SQL manuales:

```bash
make migrate-create name=nombre_de_tu_migracion
```

Esto creará dos archivos en la carpeta `migrations/`:
- `XXXXXX_nombre_de_tu_migracion.up.sql`
- `XXXXXX_nombre_de_tu_migracion.down.sql`

Para ejecutar migraciones:

```bash
make migrate-up
```

Para revertir la última migración:

```bash
make migrate-down
```

## 🔐 Seguridad

### En Producción

1. **Cambia el JWT_SECRET**: Usa una clave segura y aleatoria
2. **Cambia las contraseñas**: Especialmente la del usuario admin
3. **Usa HTTPS**: Configura un certificado SSL/TLS
4. **Configura CORS**: Limita los orígenes permitidos
5. **Variables de entorno**: No commitees el archivo `.env`

## 🐛 Troubleshooting

### Error: "connection refused"

Asegúrate de que PostgreSQL esté corriendo:

```bash
docker-compose ps
```

### Error: "database does not exist"

Crea la base de datos manualmente o reinicia los contenedores:

```bash
docker-compose down -v
docker-compose up -d
```

### Error: "port already in use"

Cambia el puerto en el archivo `.env` o detén el proceso que está usando el puerto 8080.

## 📚 Documentación de la API

Consulta el archivo `README.md` principal para ver todos los endpoints disponibles.

## 🤝 Contribuir

1. Crea una rama para tu feature: `git checkout -b feature/nueva-funcionalidad`
2. Haz commit de tus cambios: `git commit -m 'Agregar nueva funcionalidad'`
3. Push a la rama: `git push origin feature/nueva-funcionalidad`
4. Abre un Pull Request

## 📄 Licencia

Este proyecto es privado y confidencial.
