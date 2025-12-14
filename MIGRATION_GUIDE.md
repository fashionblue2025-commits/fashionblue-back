# 🗄️ Guía de Migraciones en Producción

Guía para ejecutar migraciones de base de datos en el VPS de producción.

## 🚀 Métodos de Ejecución

### Método 1: Script Automatizado (Recomendado)

El método más seguro y fácil:

```bash
# SSH al VPS
ssh root@72.60.167.46

# Ir al directorio del proyecto
cd /opt/fashion-blue

# Ejecutar migraciones
bash scripts/migrate-prod.sh
```

El script:
- ✅ Verifica que el contenedor PostgreSQL esté corriendo
- ✅ Crea tabla de control de migraciones
- ✅ Evita ejecutar migraciones duplicadas
- ✅ Muestra un resumen al final
- ✅ Registra qué migraciones se aplicaron y cuándo

---

### Método 2: Manualmente (Una por una)

Si prefieres control total:

```bash
# SSH al VPS
ssh root@72.60.167.46
cd /opt/fashion-blue

# Ver migraciones disponibles
ls -la migrations/

# Ejecutar una migración específica
docker exec -i fashionblue-postgres-prod \
  psql -U fashionblue_user -d fashionblue_prod \
  < migrations/tu_migracion.sql
```

---

### Método 3: Usando golang-migrate en el VPS

Instalar golang-migrate en el VPS y usarlo directamente:

```bash
# SSH al VPS
ssh root@72.60.167.46

# Instalar golang-migrate
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Verificar instalación
migrate -version

# Ejecutar migraciones
cd /opt/fashion-blue

# Cargar variables de entorno
source .env

# Ejecutar
migrate -path migrations \
  -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" \
  up
```

---

## 📋 Migraciones Disponibles

Actualmente tienes estas migraciones en el proyecto:

```
migrations/
├── add_category_id_to_order_items.sql
├── add_customer_id_to_orders.sql
├── add_reserved_quantity_to_order_items.sql
├── add_slug_to_categories.sql
├── create_audit_logs_table.sql
├── create_financial_transactions_table.sql
├── migrate_to_financial_transactions.sql
└── remove_quantity_reserved_from_order_items.sql
```

---

## 🔍 Verificar Estado de Migraciones

### Ver qué migraciones se han aplicado:

```bash
# SSH al VPS
ssh root@72.60.167.46

# Conectarse a PostgreSQL
docker exec -it fashionblue-postgres-prod \
  psql -U fashionblue_user -d fashionblue_prod

# Listar migraciones aplicadas
SELECT * FROM schema_migrations ORDER BY applied_at;

# Salir
\q
```

### Ver estructura de tablas:

```bash
docker exec -it fashionblue-postgres-prod \
  psql -U fashionblue_user -d fashionblue_prod

# Listar todas las tablas
\dt

# Ver estructura de una tabla específica
\d nombre_tabla

# Salir
\q
```

---

## ⚠️ Antes de Ejecutar Migraciones

### Checklist:

- [ ] **Backup creado** - Siempre crear backup antes de migraciones
- [ ] **Contenedores corriendo** - Verificar que PostgreSQL esté up
- [ ] **Revisar SQL** - Leer las migraciones que se van a aplicar
- [ ] **Ambiente correcto** - Confirmar que estás en producción

### Crear Backup:

```bash
# SSH al VPS
ssh root@72.60.167.46
cd /opt/fashion-blue

# Crear backup manual
bash scripts/backup-db.sh

# Verificar que se creó
ls -lh backups/postgres/
```

---

## 🆘 Rollback / Revertir Migraciones

Si algo sale mal:

### Opción 1: Restaurar desde backup

```bash
# SSH al VPS
ssh root@72.60.167.46
cd /opt/fashion-blue

# Listar backups disponibles
ls -lh backups/postgres/

# Restaurar
bash scripts/restore-db.sh
# Selecciona el backup que quieres restaurar
```

### Opción 2: Revertir manualmente

Si tienes migraciones down (reversibles), créalas como `*.down.sql`:

```bash
# Ejecutar migración down
docker exec -i fashionblue-postgres-prod \
  psql -U fashionblue_user -d fashionblue_prod \
  < migrations/tu_migracion.down.sql
```

---

## 🎯 Ejemplo Completo

Ejecutar todas las migraciones con el método recomendado:

```bash
# 1. Conectarse al VPS
ssh root@72.60.167.46

# 2. Ir al proyecto
cd /opt/fashion-blue

# 3. Crear backup de seguridad
echo "📦 Creando backup..."
bash scripts/backup-db.sh

# 4. Verificar contenedores
echo "🔍 Verificando contenedores..."
docker compose -f docker-compose.prod.yml ps

# 5. Ejecutar migraciones
echo "🚀 Ejecutando migraciones..."
bash scripts/migrate-prod.sh

# 6. Verificar que todo funcione
echo "✅ Verificando API..."
docker compose -f docker-compose.prod.yml logs --tail=50 api
```

---

## 📊 Monitoreo Post-Migración

Después de ejecutar migraciones:

```bash
# Ver logs de la API
docker compose -f docker-compose.prod.yml logs -f api

# Ver logs de PostgreSQL
docker compose -f docker-compose.prod.yml logs -f postgres

# Verificar que la API responde
curl http://localhost:8080/health
```

---

## 🔧 Troubleshooting

### Error: "Tabla ya existe"

Alguna migración ya fue aplicada. Revisa:

```sql
SELECT * FROM schema_migrations;
```

Si es correcto, modifica la migración para usar `CREATE TABLE IF NOT EXISTS`.

### Error: "Contenedor no está corriendo"

```bash
docker compose -f docker-compose.prod.yml up -d postgres
```

### Error: "Permission denied"

```bash
# Dar permisos al script
chmod +x scripts/migrate-prod.sh

# O ejecutar con bash explícitamente
bash scripts/migrate-prod.sh
```

### Error de conexión a base de datos

Verifica las credenciales en `.env`:

```bash
cat .env | grep -E "DB_USER|DB_PASSWORD|DB_NAME"
```

---

## 📝 Crear Nueva Migración

Para desarrollo futuro:

```bash
# En tu máquina local
cd /Users/bryanarroyaveortiz/Documents/PERSONAL/Proyectos/fashion-blue

# Crear nueva migración
migrate create -ext sql -dir migrations -seq nombre_descriptivo

# Esto creará:
# migrations/000001_nombre_descriptivo.up.sql
# migrations/000001_nombre_descriptivo.down.sql

# Editar los archivos y hacer commit
git add migrations/
git commit -m "Add migration: nombre_descriptivo"
git push origin main
```

El deployment automático subirá las migraciones al VPS, pero NO las ejecutará automáticamente por seguridad.

---

## 🎓 Best Practices

1. **Siempre hacer backup** antes de migraciones
2. **Revisar el SQL** antes de ejecutar
3. **Probar en desarrollo** primero
4. **Una migración = una acción** (no mezclar múltiples cambios)
5. **Nombrar descriptivamente** los archivos
6. **Mantener orden** (usar números de secuencia)
7. **Documentar cambios** complejos
8. **Crear migraciones down** cuando sea posible

---

## 🚀 Automatizar Migraciones (Opcional)

Si quieres que las migraciones se ejecuten automáticamente en cada deployment:

**NO RECOMENDADO** para producción, pero si lo deseas, modifica `.github/workflows/deploys.yml`:

```yaml
# Agregar después del paso "Deploy with zero-downtime"
- name: Run migrations
  run: |
    echo "🗄️ Running database migrations..."
    bash scripts/migrate-prod.sh
```

**⚠️ Cuidado:** Esto ejecutará migraciones automáticamente en cada deploy.

---

**Creado por:** Bryan Arroyave  
**Proyecto:** Fashion Blue API  
**Última actualización:** Diciembre 2025
