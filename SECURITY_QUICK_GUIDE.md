# 🔒 Guía Rápida de Seguridad

## ⚠️ IMPORTANTE: Antes de Desplegar

### 1. Configurar Variables de Entorno

```bash
# Desarrollo
cp .env.development.example .env.development
# Editar y usar credenciales simples

# Producción
cp .env.production.example .env.production
# ⚠️ CAMBIAR TODAS las credenciales marcadas con CHANGE_ME
```

### 2. Generar Secretos Seguros

```bash
# JWT Secret
openssl rand -base64 64

# Contraseñas
openssl rand -base64 32
```

### 3. Verificar que NO estén en Git

```bash
# Estos archivos NO deben estar en el repositorio:
.env
.env.development
.env.production
.env.staging
```

## 🚀 Despliegue Seguro

### Desarrollo
```bash
./scripts/deploy-dev.sh
```

### Producción
```bash
# ⚠️ Verificar TODAS las credenciales antes
./scripts/deploy-prod.sh
```

## ✅ Checklist Mínimo

- [ ] Cambié todas las contraseñas de ejemplo
- [ ] Generé un JWT_SECRET único
- [ ] Configuré CORS solo para mis dominios
- [ ] Habilité SSL en producción (DB_SSLMODE=require)
- [ ] Los archivos .env.* están en .gitignore
- [ ] pgAdmin está deshabilitado en producción

## 🆘 Si algo sale mal

1. Detener servicios: `docker-compose down`
2. Revisar logs: `docker-compose logs`
3. Ver guía completa: `SECURITY.md`

---

**Ver `SECURITY.md` para la guía completa de seguridad.**
