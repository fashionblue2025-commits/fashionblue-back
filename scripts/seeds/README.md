# 🌱 Database Seeds - Fashion Blue

Scripts SQL para poblar la base de datos con datos iniciales.

## 📋 Contenido

### Archivos SQL

1. **`01_users.sql`** - Usuario Super Admin
2. **`02_categories.sql`** - Categorías de productos (5)
3. **`03_sizes.sql`** - Tallas de ropa y calzado (33)
4. **`seed_all.sql`** - Ejecuta todos los seeds en orden

## 🚀 Cómo Ejecutar

### Opción 1: Script Bash (Recomendado)

```bash
# Dar permisos de ejecución
chmod +x scripts/run-seeds.sh

# Ejecutar
./scripts/run-seeds.sh
```

El script:
- ✅ Lee las credenciales de `.env`
- ✅ Verifica la conexión a PostgreSQL
- ✅ Ejecuta todos los seeds en orden
- ✅ Muestra un resumen al final

### Opción 2: psql Directo

```bash
# Ejecutar todos los seeds
psql -h localhost -U fashionblue -d fashionblue_db -f scripts/seeds/seed_all.sql

# O ejecutar seeds individuales
psql -h localhost -U fashionblue -d fashionblue_db -f scripts/seeds/01_users.sql
psql -h localhost -U fashionblue -d fashionblue_db -f scripts/seeds/02_categories.sql
psql -h localhost -U fashionblue -d fashionblue_db -f scripts/seeds/03_sizes.sql
```

### Opción 3: Desde Docker

```bash
# Si PostgreSQL está en Docker
docker exec -i fashionblue-postgres psql -U fashionblue -d fashionblue_db < scripts/seeds/seed_all.sql

# O ejecutar seeds individuales
docker exec -i fashionblue-postgres psql -U fashionblue -d fashionblue_db < scripts/seeds/01_users.sql
docker exec -i fashionblue-postgres psql -U fashionblue -d fashionblue_db < scripts/seeds/02_categories.sql
docker exec -i fashionblue-postgres psql -U fashionblue -d fashionblue_db < scripts/seeds/03_sizes.sql
```

### Opción 4: pgAdmin

1. Abrir pgAdmin en http://localhost:5050
2. Conectar al servidor
3. Abrir Query Tool
4. Copiar y pegar el contenido de cada archivo SQL
5. Ejecutar

## 📊 Datos que se Crean

### 👤 Usuario Super Admin

```
Email: admin@fashionblue.com
Password: admin123
Role: SUPER_ADMIN
```

### 📁 Categorías (5)

1. Chaquetas - Chaquetas y abrigos de cuero
2. Pantalones - Pantalones y jeans de cuero
3. Camisas - Camisas y blusas
4. Accesorios - Cinturones, carteras y más
5. Calzado - Zapatos y botas de cuero

### 📏 Tallas (33 total)

#### Camisetas (6)
- XS, S, M, L, XL, XXL

#### Pantalones (10)
- 24, 26, 28, 30, 32, 34, 36, 38, 40, 42 (pulgadas)

#### Zapatos (17)
- 5, 5.5, 6, 6.5, 7, 7.5, 8, 8.5, 9, 9.5, 10, 10.5, 11, 11.5, 12, 13, 14 (US)

## 🔄 Re-ejecutar Seeds

Los seeds usan `ON CONFLICT DO NOTHING`, por lo que puedes ejecutarlos múltiples veces sin crear duplicados.

```bash
# Es seguro ejecutar múltiples veces
./scripts/run-seeds.sh
```

## 🗑️ Limpiar y Re-seed

Si quieres limpiar todo y empezar de nuevo:

```bash
# Opción 1: Eliminar datos manualmente
psql -h localhost -U fashionblue -d fashionblue_db -c "TRUNCATE users, categories, sizes CASCADE;"

# Opción 2: Eliminar y recrear base de datos
docker-compose down -v
docker-compose up -d postgres
# Esperar a que PostgreSQL esté listo
go run cmd/api/main.go  # Ejecuta migraciones
./scripts/run-seeds.sh  # Ejecuta seeds
```

## ⚠️ Notas Importantes

### Hash de Contraseña

El usuario admin usa bcrypt para hashear la contraseña. El hash incluido es para `admin123`:

```
$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
```

Si necesitas cambiar la contraseña, genera un nuevo hash:

```bash
# En Go
go run -c 'package main; import ("fmt"; "golang.org/x/crypto/bcrypt"); func main() { hash, _ := bcrypt.GenerateFromPassword([]byte("tu_password"), bcrypt.DefaultCost); fmt.Println(string(hash)) }'

# O usa una herramienta online (no recomendado para producción)
# https://bcrypt-generator.com/
```

### Orden de Ejecución

Los seeds deben ejecutarse en orden:
1. Primero `01_users.sql` (no tiene dependencias)
2. Luego `02_categories.sql` (no tiene dependencias)
3. Finalmente `03_sizes.sql` (no tiene dependencias)

Actualmente no hay dependencias entre tablas en los seeds, pero es buena práctica mantener el orden.

## 🧪 Verificar Seeds

```bash
# Verificar usuario
psql -h localhost -U fashionblue -d fashionblue_db -c "SELECT email, first_name, last_name, role FROM users;"

# Verificar categorías
psql -h localhost -U fashionblue -d fashionblue_db -c "SELECT name, description FROM categories;"

# Verificar tallas
psql -h localhost -U fashionblue -d fashionblue_db -c "SELECT type, COUNT(*) FROM sizes GROUP BY type;"

# Contar todos los registros
psql -h localhost -U fashionblue -d fashionblue_db -c "
SELECT 
    (SELECT COUNT(*) FROM users) as users,
    (SELECT COUNT(*) FROM categories) as categories,
    (SELECT COUNT(*) FROM sizes) as sizes;
"
```

## 🔐 Seguridad

**⚠️ IMPORTANTE:** 

- El usuario `admin@fashionblue.com` con password `admin123` es solo para desarrollo
- **NUNCA uses estas credenciales en producción**
- Cambia la contraseña inmediatamente después del primer login en producción
- Usa variables de entorno para credenciales sensibles

## 📚 Más Información

- Ver `API_EXAMPLES.md` para ejemplos de uso de la API
- Ver `CUSTOMER_REFACTOR_GUIDE.md` para información sobre el modelo de clientes
- Ver `SECURITY.md` para mejores prácticas de seguridad

---

**Última actualización:** 2024-11-20
