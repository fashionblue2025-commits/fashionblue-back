# Fashion Blue - Sistema de Gestión Empresarial

Sistema completo de gestión para empresa de manufactura y venta de ropa, desarrollado con arquitectura hexagonal en Go.

## 🏗️ Arquitectura

El proyecto sigue una **arquitectura hexagonal (ports & adapters)** con las siguientes capas:

```
cmd/                    # Punto de entrada de la aplicación
internal/
  ├── domain/          # Entidades y lógica de negocio
  ├── application/     # Casos de uso
  ├── ports/           # Interfaces (input/output ports)
  └── adapters/        # Implementaciones de los ports
      ├── http/        # Handlers HTTP (Echo)
      ├── postgres/    # Repositorios PostgreSQL
      └── middleware/  # Middlewares
pkg/                   # Paquetes compartidos
migrations/            # Migraciones de base de datos
```

## 🚀 Características

### Gestión Financiera
- **Inyección de Capital**: Registro y seguimiento de inversiones
- **Contabilidad**: Control de ingresos, gastos y ganancias

### Gestión de Productos
- Categorías de productos (chaquetas, pantalones, etc.)
- Costos de producción (materiales + mano de obra)
- Precios: unitario y por mayor
- Inventario y stock

### Gestión de Ventas
- Registro de ventas por vendedor
- Seguimiento de productos vendidos
- Cálculo de ganancias

### Gestión de Proveedores
- Información de proveedores
- Historial de compras
- Adjuntar facturas

### Gestión de Clientes
- Información de clientes
- Historial de movimientos y compras

### Roles y Permisos
- **Super Admin**: Control total del sistema
- **Vendedor**: Gestión de ventas
- Autenticación y autorización

## 🛠️ Tecnologías

- **Go 1.21+**
- **Echo Framework**: Framework HTTP
- **PostgreSQL**: Base de datos
- **Docker & Docker Compose**: Containerización
- **GORM**: ORM para Go
- **JWT**: Autenticación
- **golang-migrate**: Migraciones de BD

## 📦 Instalación

### Prerrequisitos
- Go 1.21 o superior
- Docker y Docker Compose
- Make (opcional)

### Configuración

1. Clonar el repositorio:
```bash
git clone <repository-url>
cd fashion-blue
```

2. Copiar el archivo de configuración:
```bash
cp .env.example .env
```

3. Configurar variables de entorno en `.env`

4. Levantar los servicios con Docker:
```bash
docker-compose up -d
```

5. Ejecutar migraciones:
```bash
make migrate-up
```

6. Iniciar la aplicación:
```bash
make run
```

## 🔧 Comandos Útiles

```bash
# Ejecutar la aplicación
make run

# Ejecutar tests
make test

# Ejecutar migraciones
make migrate-up
make migrate-down

# Limpiar y reconstruir
make clean
make build

# Ver logs de Docker
docker-compose logs -f
```

## 📝 API Endpoints

### Autenticación
- `POST /api/v1/auth/login` - Iniciar sesión
- `POST /api/v1/auth/register` - Registrar usuario

### Inyección de Capital
- `POST /api/v1/capital-injections` - Registrar inyección
- `GET /api/v1/capital-injections` - Listar inyecciones
- `GET /api/v1/capital-injections/:id` - Obtener detalle

### Productos
- `POST /api/v1/products` - Crear producto
- `GET /api/v1/products` - Listar productos
- `GET /api/v1/products/:id` - Obtener producto
- `PUT /api/v1/products/:id` - Actualizar producto
- `DELETE /api/v1/products/:id` - Eliminar producto

### Categorías
- `POST /api/v1/categories` - Crear categoría
- `GET /api/v1/categories` - Listar categorías

### Ventas
- `POST /api/v1/sales` - Registrar venta
- `GET /api/v1/sales` - Listar ventas
- `GET /api/v1/sales/:id` - Obtener venta
- `GET /api/v1/sales/stats` - Estadísticas de ventas

### Proveedores
- `POST /api/v1/suppliers` - Crear proveedor
- `GET /api/v1/suppliers` - Listar proveedores
- `GET /api/v1/suppliers/:id` - Obtener proveedor
- `PUT /api/v1/suppliers/:id` - Actualizar proveedor

### Compras a Proveedores
- `POST /api/v1/purchases` - Registrar compra
- `GET /api/v1/purchases` - Listar compras
- `POST /api/v1/purchases/:id/invoice` - Adjuntar factura

### Clientes
- `POST /api/v1/customers` - Crear cliente
- `GET /api/v1/customers` - Listar clientes
- `GET /api/v1/customers/:id` - Obtener cliente
- `GET /api/v1/customers/:id/history` - Historial del cliente

### Usuarios (Super Admin)
- `POST /api/v1/users` - Crear usuario
- `GET /api/v1/users` - Listar usuarios
- `PUT /api/v1/users/:id` - Actualizar usuario
- `DELETE /api/v1/users/:id` - Eliminar usuario

## 🔐 Roles y Permisos

- **SUPER_ADMIN**: Acceso total al sistema
- **SELLER**: Gestión de ventas y clientes

## 📊 Base de Datos

El sistema utiliza PostgreSQL con las siguientes tablas principales:

- `users` - Usuarios del sistema
- `capital_injections` - Inyecciones de capital
- `categories` - Categorías de productos
- `products` - Productos
- `sales` - Ventas
- `sale_items` - Ítems de venta
- `suppliers` - Proveedores
- `purchases` - Compras a proveedores
- `purchase_items` - Ítems de compra
- `customers` - Clientes
- `customer_transactions` - Transacciones de clientes

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto es privado y confidencial.

## 👥 Equipo

Fashion Blue Development Team
