# 🔒 Guía de Seguridad - Fashion Blue

## 📋 Índice
1. [Configuración de Ambientes](#configuración-de-ambientes)
2. [Gestión de Secretos](#gestión-de-secretos)
3. [Mejores Prácticas](#mejores-prácticas)
4. [Checklist de Seguridad](#checklist-de-seguridad)
5. [Respuesta a Incidentes](#respuesta-a-incidentes)

---

## 🌍 Configuración de Ambientes

### Desarrollo Local

```bash
# 1. Copiar archivo de ejemplo
cp .env.development.example .env.development

# 2. Editar con valores de desarrollo
nano .env.development

# 3. Levantar servicios
./scripts/deploy-dev.sh
```

**Características:**
- ✅ Credenciales simples (no críticas)
- ✅ Logs verbosos
- ✅ CORS permisivo
- ✅ pgAdmin habilitado

### Staging

```bash
# 1. Copiar archivo de ejemplo
cp .env.staging.example .env.staging

# 2. Configurar con credenciales de staging
nano .env.staging

# 3. Desplegar
docker-compose --env-file .env.staging up -d
```

**Características:**
- ⚠️ Credenciales intermedias
- ⚠️ SSL habilitado
- ⚠️ CORS restrictivo
- ⚠️ pgAdmin opcional

### Producción

```bash
# 1. Copiar archivo de ejemplo
cp .env.production.example .env.production

# 2. Configurar con credenciales FUERTES
nano .env.production

# 3. Verificar configuración
./scripts/deploy-prod.sh
```

**Características:**
- 🔒 Credenciales fuertes y únicas
- 🔒 SSL obligatorio
- 🔒 CORS muy restrictivo
- 🔒 pgAdmin deshabilitado
- 🔒 Logs en JSON
- 🔒 Backups automáticos

---

## 🔐 Gestión de Secretos

### ❌ NUNCA hacer esto:

```yaml
# ❌ MAL - Credenciales hardcodeadas
environment:
  DB_PASSWORD: mypassword123
  JWT_SECRET: supersecret
```

### ✅ SIEMPRE hacer esto:

```yaml
# ✅ BIEN - Variables de entorno
environment:
  DB_PASSWORD: ${DB_PASSWORD}
  JWT_SECRET: ${JWT_SECRET}
```

### Generar Secretos Seguros

```bash
# JWT Secret (64 caracteres)
openssl rand -base64 64

# Contraseña fuerte
openssl rand -base64 32

# UUID único
uuidgen
```

### Servicios de Gestión de Secretos (Recomendado para Producción)

#### AWS Secrets Manager
```bash
# Guardar secreto
aws secretsmanager create-secret \
  --name fashionblue/prod/db-password \
  --secret-string "your-strong-password"

# Recuperar secreto
aws secretsmanager get-secret-value \
  --secret-id fashionblue/prod/db-password \
  --query SecretString \
  --output text
```

#### HashiCorp Vault
```bash
# Guardar secreto
vault kv put secret/fashionblue/prod \
  db_password="your-strong-password" \
  jwt_secret="your-jwt-secret"

# Recuperar secreto
vault kv get -field=db_password secret/fashionblue/prod
```

#### Docker Secrets (Docker Swarm)
```bash
# Crear secret
echo "your-strong-password" | docker secret create db_password -

# Usar en docker-compose
services:
  api:
    secrets:
      - db_password
    environment:
      DB_PASSWORD_FILE: /run/secrets/db_password
```

---

## 🛡️ Mejores Prácticas

### 1. Contraseñas

✅ **Hacer:**
- Mínimo 16 caracteres
- Mezcla de mayúsculas, minúsculas, números y símbolos
- Única para cada servicio
- Rotar cada 90 días
- Usar un gestor de contraseñas

❌ **No hacer:**
- Usar contraseñas comunes (admin123, password, etc.)
- Reutilizar contraseñas
- Compartir contraseñas por email/chat
- Guardar en archivos de texto plano

### 2. JWT

```bash
# Generar clave segura
JWT_SECRET=$(openssl rand -base64 64)

# Configurar expiración corta
JWT_EXPIRATION=1h  # Para APIs sensibles
JWT_EXPIRATION=24h # Para aplicaciones normales
```

**Recomendaciones:**
- Mínimo 32 caracteres
- Rotar cada 6 meses
- Usar algoritmo HS256 o RS256
- Implementar refresh tokens
- Blacklist de tokens revocados

### 3. Base de Datos

```env
# Producción - SIEMPRE usar SSL
DB_SSLMODE=require

# Desarrollo - Solo si es necesario
DB_SSLMODE=disable
```

**Configuración segura:**
```sql
-- Crear usuario con permisos limitados
CREATE USER fashionblue_app WITH PASSWORD 'strong-password';
GRANT CONNECT ON DATABASE fashionblue_db TO fashionblue_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO fashionblue_app;

-- No dar permisos de DROP, CREATE, ALTER
```

### 4. CORS

```env
# ❌ Desarrollo - Permisivo
CORS_ALLOWED_ORIGINS=*

# ✅ Producción - Restrictivo
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

### 5. Logs

```env
# Desarrollo
LOG_LEVEL=debug
LOG_FORMAT=text

# Producción
LOG_LEVEL=info  # o warning
LOG_FORMAT=json
```

**⚠️ Nunca loguear:**
- Contraseñas
- Tokens JWT
- Números de tarjetas
- Información personal sensible

### 6. Rate Limiting

Implementar en el código:
```go
// Ejemplo con Echo
e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
```

### 7. HTTPS

**Producción - OBLIGATORIO:**
```nginx
# nginx.conf
server {
    listen 443 ssl http2;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    # Redirigir HTTP a HTTPS
    if ($scheme != "https") {
        return 301 https://$server_name$request_uri;
    }
}
```

---

## ✅ Checklist de Seguridad

### Antes de Desplegar a Producción

- [ ] Todas las contraseñas son fuertes y únicas
- [ ] JWT_SECRET tiene mínimo 32 caracteres aleatorios
- [ ] DB_SSLMODE está en `require`
- [ ] CORS está configurado solo para dominios permitidos
- [ ] No hay credenciales hardcodeadas en el código
- [ ] `.env.production` está en `.gitignore`
- [ ] Logs no contienen información sensible
- [ ] HTTPS está habilitado
- [ ] Rate limiting está implementado
- [ ] Backups automáticos están configurados
- [ ] Monitoreo y alertas están activos
- [ ] pgAdmin está deshabilitado o protegido
- [ ] Firewall está configurado
- [ ] Actualizaciones de seguridad están aplicadas

### Revisión Periódica (Cada 3 meses)

- [ ] Rotar credenciales
- [ ] Revisar logs de acceso
- [ ] Actualizar dependencias
- [ ] Revisar permisos de usuarios
- [ ] Verificar backups
- [ ] Auditar código
- [ ] Penetration testing

---

## 🚨 Respuesta a Incidentes

### Si se compromete una credencial:

1. **Inmediato (< 5 minutos):**
   ```bash
   # Cambiar credencial comprometida
   # Reiniciar servicios afectados
   docker-compose restart api
   ```

2. **Corto plazo (< 1 hora):**
   - Revisar logs de acceso
   - Identificar accesos no autorizados
   - Revocar tokens activos
   - Notificar al equipo

3. **Mediano plazo (< 24 horas):**
   - Investigar causa raíz
   - Implementar medidas preventivas
   - Documentar incidente
   - Actualizar procedimientos

4. **Largo plazo (< 1 semana):**
   - Auditoría completa de seguridad
   - Capacitación del equipo
   - Mejorar monitoreo
   - Post-mortem

### Contactos de Emergencia

```
Security Lead: security@yourdomain.com
DevOps Lead:   devops@yourdomain.com
On-Call:       +1-XXX-XXX-XXXX
```

---

## 📚 Recursos Adicionales

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [CIS Docker Benchmark](https://www.cisecurity.org/benchmark/docker)
- [Go Security Checklist](https://github.com/Checkmarx/Go-SCP)
- [PostgreSQL Security](https://www.postgresql.org/docs/current/security.html)

---

## 🔄 Historial de Cambios

| Fecha | Cambio | Responsable |
|-------|--------|-------------|
| 2024-11-20 | Documento inicial | DevOps Team |

---

**⚠️ Este documento debe revisarse y actualizarse regularmente.**

**Última actualización:** 2024-11-20
