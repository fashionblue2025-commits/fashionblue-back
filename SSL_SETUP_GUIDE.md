# 🔒 Guía de Configuración SSL/HTTPS

Guía paso a paso para configurar HTTPS en Fashion Blue API usando dominio + Let's Encrypt SSL.

---

## 📋 Checklist General

- [ ] Comprar dominio
- [ ] Configurar DNS
- [ ] Esperar propagación DNS
- [ ] Ejecutar script de SSL
- [ ] Actualizar frontend
- [ ] Verificar funcionamiento

---

## 🛒 Paso 1: Comprar Dominio

### Opción Recomendada: Namecheap

1. **Ir a:** https://www.namecheap.com
2. **Buscar dominio:** `fashionblue` (o el que prefieras)
3. **Seleccionar extensión:**
   - `.com` → $8.88/año (profesional)
   - `.tech` → $2.98/año (moderno)
   - `.xyz` → $0.98/año primer año (económico)

4. **Sugerencias de nombres:**
   - `fashionblue-api.com`
   - `api.fashionblue.tech`
   - `fashionblue.tech`
   - `fbapi.com`

5. **Completar compra** y crear cuenta

---

## 🌐 Paso 2: Configurar DNS

Una vez comprado el dominio, configurar el DNS:

### En Namecheap:

1. **Ir a:** Dashboard → Domain List
2. **Clic en:** Manage (al lado de tu dominio)
3. **Ir a:** Advanced DNS
4. **Agregar registro A:**

```
Type: A Record
Host: @ (para dominio raíz) o "api" (para subdominio)
Value: 72.60.167.46
TTL: Automatic (o 300)
```

**Ejemplos:**

Si quieres `fashionblue.com`:
```
Host: @
Value: 72.60.167.46
```

Si quieres `api.fashionblue.com`:
```
Host: api
Value: 72.60.167.46
```

5. **Guardar cambios**

### En Cloudflare (si lo usas):

1. **Agregar sitio** a Cloudflare
2. **Crear registro A:**
```
Type: A
Name: @ o api
IPv4: 72.60.167.46
Proxy status: DNS only (nube gris, NO naranja)
```

**⚠️ Importante:** Desactiva el proxy (nube gris) para Let's Encrypt

---

## ⏳ Paso 3: Esperar Propagación DNS

Verificar que el dominio apunte correctamente:

```bash
# En tu Mac, verificar DNS
nslookup tudominio.com

# O con dig
dig tudominio.com

# Debe mostrar: 72.60.167.46
```

**Tiempo de espera:** 5 minutos a 48 horas (usualmente 10-30 minutos)

**Verificar online:**
- https://dnschecker.org/

---

## 🚀 Paso 4: Ejecutar Script de Configuración SSL

### 4.1 Subir script al repositorio

```bash
# En tu Mac
cd /Users/bryanarroyaveortiz/Documents/PERSONAL/Proyectos/fashion-blue

git add scripts/setup-ssl.sh
git commit -m "Add SSL setup script"
git push origin main
```

### 4.2 En el VPS, ejecutar el script

```bash
# SSH al VPS
ssh root@72.60.167.46

# Ir al proyecto
cd /opt/fashion-blue

# Pull del script
git pull origin main

# Dar permisos de ejecución
chmod +x scripts/setup-ssl.sh

# Ejecutar script
bash scripts/setup-ssl.sh
```

### 4.3 Seguir las instrucciones del script

El script te pedirá:
1. **Dominio:** Ingresa tu dominio completo (ej: `api.fashionblue.com`)
2. **Email:** Tu email para notificaciones de SSL

El script hará automáticamente:
- ✅ Instalar Nginx
- ✅ Configurar Nginx como reverse proxy
- ✅ Instalar Certbot
- ✅ Obtener certificado SSL de Let's Encrypt
- ✅ Configurar renovación automática
- ✅ Configurar firewall
- ✅ Verificar instalación

---

## 🔄 Paso 5: Actualizar Frontend

### En el código del frontend:

```javascript
// Antes (HTTP)
const API_BASE_URL = 'http://72.60.167.46';

// Después (HTTPS)
const API_BASE_URL = 'https://tudominio.com';
// o
const API_BASE_URL = 'https://api.tudominio.com';
```

### Actualizar variable de entorno:

```bash
# .env o .env.production en el frontend
VITE_API_URL=https://tudominio.com
# o
REACT_APP_API_URL=https://tudominio.com
# o
NEXT_PUBLIC_API_URL=https://tudominio.com
```

### Rebuild y redeploy del frontend:

```bash
npm run build
# Subir a Netlify/Vercel
```

---

## ✅ Paso 6: Verificar Funcionamiento

### 6.1 Verificar SSL

En tu navegador, visita:
```
https://tudominio.com/health
```

Debe mostrar:
- 🔒 Candado verde en la barra de direcciones
- Respuesta JSON del health check

### 6.2 Verificar certificado

```bash
# Ver detalles del certificado
sudo certbot certificates

# Probar renovación (dry run)
sudo certbot renew --dry-run
```

### 6.3 Test desde frontend

1. Ir a tu frontend en Netlify
2. Intentar login/registro
3. **NO** debe aparecer el error de "Mixed Content"
4. Las peticiones deben funcionar correctamente

---

## 🔧 Troubleshooting

### Error: "DNS no apunta al servidor"

**Solución:**
- Espera más tiempo (propagación DNS)
- Verifica que configuraste el registro A correctamente
- Usa `nslookup tudominio.com` para verificar

### Error: "Puerto 80 o 443 en uso"

**Solución:**
```bash
# Ver qué está usando el puerto
sudo lsof -i :80
sudo lsof -i :443

# Detener servicios conflictivos
sudo systemctl stop apache2  # Si tienes Apache
```

### Error: "Certificado no válido"

**Solución:**
```bash
# Forzar renovación
sudo certbot renew --force-renewal
```

### Error: Frontend sigue con "Mixed Content"

**Solución:**
- Verifica que actualizaste la variable de entorno
- Haz rebuild del frontend
- Limpia cache del navegador (Ctrl + Shift + R)
- Verifica en las DevTools → Network que las peticiones van a HTTPS

---

## 🔄 Renovación Automática

El certificado SSL se renueva automáticamente cada 90 días.

### Verificar servicio de renovación:

```bash
# Ver estado del timer
sudo systemctl status certbot.timer

# Ver logs de renovación
sudo journalctl -u certbot.renew
```

### Renovar manualmente (si necesario):

```bash
sudo certbot renew
sudo systemctl reload nginx
```

---

## 📊 Resumen de URLs

Una vez configurado:

| Tipo | URL Anterior | URL Nueva |
|------|--------------|-----------|
| **API Base** | `http://72.60.167.46` | `https://tudominio.com` |
| **Health Check** | `http://72.60.167.46/health` | `https://tudominio.com/health` |
| **Login** | `http://72.60.167.46/api/v1/auth/login` | `https://tudominio.com/api/v1/auth/login` |

---

## 🎯 Checklist Final

Después de configurar, verifica:

- [ ] `https://tudominio.com/health` responde correctamente
- [ ] Candado verde en el navegador
- [ ] Frontend puede hacer peticiones sin error de "Mixed Content"
- [ ] Login/registro funciona desde el frontend
- [ ] Renovación automática está habilitada

---

## 📞 Soporte

Si tienes problemas:

1. **Ver logs de Nginx:**
   ```bash
   sudo tail -f /var/log/nginx/fashionblue-error.log
   ```

2. **Ver logs de Certbot:**
   ```bash
   sudo cat /var/log/letsencrypt/letsencrypt.log
   ```

3. **Ver estado de servicios:**
   ```bash
   sudo systemctl status nginx
   sudo systemctl status docker
   ```

4. **Verificar contenedores Docker:**
   ```bash
   docker compose -f docker-compose.prod.yml ps
   docker compose -f docker-compose.prod.yml logs -f api
   ```

---

**Creado por:** Bryan Arroyave  
**Proyecto:** Fashion Blue API  
**Última actualización:** Diciembre 2025
