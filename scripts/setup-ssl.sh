#!/bin/bash

###############################################################################
# Script para configurar HTTPS/SSL con Let's Encrypt
# Prerrequisito: Dominio apuntando a la IP del VPS (72.60.167.46)
###############################################################################

set -e

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

clear
echo -e "${BLUE}"
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║         Fashion Blue - HTTPS/SSL Configuration                ║"
echo "║              Let's Encrypt + Nginx                             ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo -e "${NC}\n"

# Verificar que estamos en el VPS
if [ ! -d "/opt/fashion-blue" ]; then
    echo -e "${RED}❌ Este script debe ejecutarse en el VPS${NC}"
    exit 1
fi

# Pedir dominio
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Paso 1: Configuración del Dominio${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

echo -e "${YELLOW}Ingresa tu dominio completo (ejemplo: api.fashionblue.com):${NC}"
read -r DOMAIN

if [ -z "$DOMAIN" ]; then
    echo -e "${RED}❌ Dominio es requerido${NC}"
    exit 1
fi

echo -e "${YELLOW}Ingresa tu email para las notificaciones de SSL:${NC}"
read -r EMAIL

if [ -z "$EMAIL" ]; then
    echo -e "${RED}❌ Email es requerido${NC}"
    exit 1
fi

echo -e "\n${GREEN}✓${NC} Dominio: ${CYAN}${DOMAIN}${NC}"
echo -e "${GREEN}✓${NC} Email: ${CYAN}${EMAIL}${NC}\n"

# Verificar que el dominio apunte a este servidor
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Paso 2: Verificando DNS${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

SERVER_IP=$(curl -s ifconfig.me)
DOMAIN_IP=$(dig +short ${DOMAIN} | tail -n1)

echo -e "${CYAN}IP del servidor:${NC} ${SERVER_IP}"
echo -e "${CYAN}IP del dominio:${NC} ${DOMAIN_IP}\n"

if [ "$SERVER_IP" != "$DOMAIN_IP" ]; then
    echo -e "${RED}❌ El dominio no apunta a este servidor${NC}"
    echo -e "${YELLOW}Configura un registro A en tu DNS:${NC}"
    echo -e "   Tipo: A"
    echo -e "   Nombre: @ (o el subdominio)"
    echo -e "   Valor: ${SERVER_IP}"
    echo -e "   TTL: 300\n"
    echo -e "${YELLOW}¿Deseas continuar de todas formas? (yes/no):${NC}"
    read -r continue_anyway
    if [ "$continue_anyway" != "yes" ]; then
        exit 1
    fi
else
    echo -e "${GREEN}✅ DNS configurado correctamente${NC}\n"
fi

# Actualizar sistema e instalar dependencias
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Paso 3: Instalando dependencias${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

sudo apt update
sudo apt install -y nginx certbot python3-certbot-nginx

echo -e "${GREEN}✅ Dependencias instaladas${NC}\n"

# Configurar Nginx
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Paso 4: Configurando Nginx${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

# Crear configuración de Nginx
sudo tee /etc/nginx/sites-available/fashionblue > /dev/null <<EOF
server {
    listen 80;
    listen [::]:80;
    server_name ${DOMAIN};

    # Logs
    access_log /var/log/nginx/fashionblue-access.log;
    error_log /var/log/nginx/fashionblue-error.log;

    # Proxy to Docker container
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        
        # Headers
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check endpoint (sin logs)
    location /health {
        proxy_pass http://127.0.0.1:8080/health;
        access_log off;
    }
}
EOF

# Habilitar el sitio
sudo ln -sf /etc/nginx/sites-available/fashionblue /etc/nginx/sites-enabled/

# Remover sitio default si existe
sudo rm -f /etc/nginx/sites-enabled/default

# Verificar configuración
echo -e "${CYAN}Verificando configuración de Nginx...${NC}"
if sudo nginx -t; then
    echo -e "${GREEN}✅ Configuración de Nginx válida${NC}\n"
else
    echo -e "${RED}❌ Error en la configuración de Nginx${NC}"
    exit 1
fi

# Reiniciar Nginx
sudo systemctl restart nginx
sudo systemctl enable nginx

echo -e "${GREEN}✅ Nginx configurado y corriendo${NC}\n"

# Configurar firewall
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Paso 5: Configurando Firewall${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

sudo ufw allow 'Nginx Full'
sudo ufw allow 443/tcp

echo -e "${GREEN}✅ Firewall configurado${NC}\n"

# Obtener certificado SSL
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Paso 6: Obteniendo certificado SSL${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

echo -e "${YELLOW}Esto puede tardar unos segundos...${NC}\n"

sudo certbot --nginx \
    -d ${DOMAIN} \
    --non-interactive \
    --agree-tos \
    --email ${EMAIL} \
    --redirect

if [ $? -eq 0 ]; then
    echo -e "\n${GREEN}✅ Certificado SSL obtenido exitosamente${NC}\n"
else
    echo -e "\n${RED}❌ Error al obtener certificado SSL${NC}"
    echo -e "${YELLOW}Verifica que:${NC}"
    echo -e "  1. El dominio apunte correctamente a este servidor"
    echo -e "  2. Los puertos 80 y 443 estén abiertos"
    echo -e "  3. No haya otro servicio usando el puerto 80"
    exit 1
fi

# Configurar renovación automática
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Paso 7: Configurando renovación automática${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

sudo systemctl enable certbot.timer
sudo systemctl start certbot.timer

echo -e "${GREEN}✅ Renovación automática configurada${NC}\n"

# Verificar que todo funcione
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Paso 8: Verificando instalación${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

echo -e "${CYAN}Probando HTTPS...${NC}"
sleep 2

if curl -f -s -o /dev/null https://${DOMAIN}/health; then
    echo -e "${GREEN}✅ HTTPS funcionando correctamente${NC}\n"
else
    echo -e "${YELLOW}⚠️  El endpoint /health no responde (esto es normal si la app no está corriendo)${NC}\n"
fi

# Resumen final
echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                   🎉 CONFIGURACIÓN COMPLETADA 🎉                ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}\n"

echo -e "${GREEN}✨ Tu API ahora está disponible en HTTPS:${NC}\n"
echo -e "   ${CYAN}URL Principal:${NC} ${GREEN}https://${DOMAIN}${NC}"
echo -e "   ${CYAN}Health Check:${NC} ${GREEN}https://${DOMAIN}/health${NC}\n"

echo -e "${YELLOW}📋 Siguiente paso:${NC}"
echo -e "   Actualiza tu frontend para usar: ${GREEN}https://${DOMAIN}${NC}\n"

echo -e "${CYAN}🔒 Información del Certificado:${NC}"
sudo certbot certificates

echo -e "\n${YELLOW}ℹ️  El certificado se renovará automáticamente cada 90 días${NC}"
echo -e "${YELLOW}ℹ️  Puedes probar la renovación con: sudo certbot renew --dry-run${NC}\n"
