#!/bin/bash

###############################################################################
# Script para inicializar SSL/HTTPS con Docker Compose + Let's Encrypt
# Este script obtiene el certificado SSL por primera vez
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
echo "║       Fashion Blue - SSL Initialization (Docker Compose)      ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo -e "${NC}\n"

# Verificar que estamos en el directorio correcto
if [ ! -f "docker-compose.prod.yml" ]; then
    echo -e "${RED}❌ Archivo docker-compose.prod.yml no encontrado${NC}"
    echo -e "${YELLOW}Ejecuta este script desde /opt/fashion-blue${NC}"
    exit 1
fi

# Cargar variables de entorno
if [ ! -f ".env" ]; then
    echo -e "${RED}❌ Archivo .env no encontrado${NC}"
    exit 1
fi

source .env

# Configuración
DOMAIN="api.fashionblue.org"
EMAIL="${SSL_EMAIL:-admin@fashionblue.org}"
STAGING=${STAGING:-0}  # 0 = producción, 1 = staging (para testing)

echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Configuración${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
echo -e "  ${CYAN}Dominio:${NC} ${GREEN}${DOMAIN}${NC}"
echo -e "  ${CYAN}Email:${NC} ${GREEN}${EMAIL}${NC}"
echo -e "  ${CYAN}Modo:${NC} ${GREEN}$([ $STAGING -eq 1 ] && echo 'Staging (Test)' || echo 'Producción')${NC}\n"

# Verificar DNS
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Verificando DNS${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

SERVER_IP=$(curl -s ifconfig.me)
DOMAIN_IP=$(dig +short ${DOMAIN} | tail -n1)

echo -e "  ${CYAN}IP del servidor:${NC} ${SERVER_IP}"
echo -e "  ${CYAN}IP del dominio:${NC} ${DOMAIN_IP}\n"

if [ "$SERVER_IP" != "$DOMAIN_IP" ]; then
    echo -e "${YELLOW}⚠️  El dominio no apunta a este servidor${NC}"
    echo -e "${YELLOW}   El certificado podría fallar${NC}\n"
    echo -e "${CYAN}¿Deseas continuar de todas formas? (yes/no):${NC}"
    read -r continue_anyway
    if [ "$continue_anyway" != "yes" ]; then
        exit 1
    fi
else
    echo -e "${GREEN}✅ DNS configurado correctamente${NC}\n"
fi

# Crear directorios necesarios
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Creando directorios${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

mkdir -p certbot/conf
mkdir -p certbot/www
mkdir -p logs

echo -e "${GREEN}✅ Directorios creados${NC}\n"

# Crear configuración temporal de Nginx (sin SSL)
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Configurando Nginx temporal (HTTP)${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

# Backup de configuración HTTPS
if [ -f "nginx/conf.d/fashionblue.conf" ]; then
    cp nginx/conf.d/fashionblue.conf nginx/conf.d/fashionblue.conf.backup
fi

# Crear configuración temporal solo HTTP
cat > nginx/conf.d/fashionblue.conf.temp <<EOF
upstream api_backend {
    server api:8080;
}

server {
    listen 80;
    listen [::]:80;
    server_name ${DOMAIN};

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        proxy_pass http://api_backend;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

mv nginx/conf.d/fashionblue.conf.temp nginx/conf.d/fashionblue.conf

echo -e "${GREEN}✅ Configuración temporal lista${NC}\n"

# Levantar servicios (sin certbot aún)
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Iniciando servicios${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

docker compose -f docker-compose.prod.yml up -d postgres api nginx

echo -e "${GREEN}✅ Servicios iniciados${NC}\n"
echo -e "${YELLOW}Esperando 10 segundos para que los servicios estén listos...${NC}"
sleep 10

# Verificar que Nginx responde
echo -e "\n${CYAN}Verificando Nginx...${NC}"
if curl -f -s -o /dev/null http://localhost; then
    echo -e "${GREEN}✅ Nginx funcionando${NC}\n"
else
    echo -e "${RED}❌ Nginx no responde${NC}"
    docker compose -f docker-compose.prod.yml logs nginx
    exit 1
fi

# Obtener certificado SSL
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Obteniendo certificado SSL${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

# Argumentos para certbot
if [ $STAGING -eq 1 ]; then
    STAGING_ARG="--staging"
    echo -e "${YELLOW}⚠️  Usando modo STAGING (certificado de prueba)${NC}\n"
else
    STAGING_ARG=""
fi

# Ejecutar certbot
docker compose -f docker-compose.prod.yml run --rm certbot \
    certonly --webroot \
    --webroot-path=/var/www/certbot \
    --email ${EMAIL} \
    --agree-tos \
    --no-eff-email \
    ${STAGING_ARG} \
    -d ${DOMAIN}

if [ $? -eq 0 ]; then
    echo -e "\n${GREEN}✅ Certificado SSL obtenido exitosamente${NC}\n"
else
    echo -e "\n${RED}❌ Error al obtener certificado SSL${NC}"
    echo -e "${YELLOW}Verifica que:${NC}"
    echo -e "  1. El dominio ${DOMAIN} apunte correctamente"
    echo -e "  2. Los puertos 80 y 443 estén abiertos"
    echo -e "  3. No hay otro servicio usando estos puertos\n"
    
    # Restaurar configuración original
    if [ -f "nginx/conf.d/fashionblue.conf.backup" ]; then
        mv nginx/conf.d/fashionblue.conf.backup nginx/conf.d/fashionblue.conf
    fi
    
    exit 1
fi

# Restaurar configuración HTTPS
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Activando configuración HTTPS${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

if [ -f "nginx/conf.d/fashionblue.conf.backup" ]; then
    mv nginx/conf.d/fashionblue.conf.backup nginx/conf.d/fashionblue.conf
fi

# Reiniciar Nginx con configuración HTTPS
docker compose -f docker-compose.prod.yml restart nginx

echo -e "${GREEN}✅ Nginx reconfigurado para HTTPS${NC}\n"

# Iniciar certbot para renovaciones automáticas
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Iniciando servicio de renovación automática${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

docker compose -f docker-compose.prod.yml up -d certbot

echo -e "${GREEN}✅ Renovación automática configurada${NC}\n"

# Verificación final
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Verificación final${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

echo -e "${CYAN}Esperando 5 segundos...${NC}"
sleep 5

echo -e "\n${CYAN}Probando HTTPS...${NC}"
if curl -f -s -o /dev/null https://${DOMAIN}/health; then
    echo -e "${GREEN}✅ HTTPS funcionando correctamente${NC}\n"
else
    echo -e "${YELLOW}⚠️  HTTPS no responde aún (puede tomar unos minutos)${NC}\n"
fi

# Resumen final
echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                🎉 SSL CONFIGURADO EXITOSAMENTE 🎉               ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}\n"

echo -e "${GREEN}✨ Tu API ahora está disponible en HTTPS:${NC}\n"
echo -e "   ${CYAN}URL:${NC} ${GREEN}https://${DOMAIN}${NC}"
echo -e "   ${CYAN}Health Check:${NC} ${GREEN}https://${DOMAIN}/health${NC}\n"

echo -e "${YELLOW}📋 Siguiente paso:${NC}"
echo -e "   Actualiza tu frontend para usar: ${GREEN}https://${DOMAIN}${NC}\n"

echo -e "${CYAN}🔄 Renovación automática:${NC}"
echo -e "   El certificado se renovará automáticamente cada 12 horas\n"

echo -e "${CYAN}📊 Ver logs:${NC}"
echo -e "   docker compose -f docker-compose.prod.yml logs -f nginx\n"

# Mostrar estado de contenedores
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Estado de contenedores:${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
docker compose -f docker-compose.prod.yml ps

echo ""
