#!/bin/bash

# Script para verificar que todo esté listo para deployment con GitHub Actions

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Contadores
PASSED=0
FAILED=0
WARNINGS=0

echo -e "${BLUE}════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  Fashion Blue - Deployment Readiness Check    ${NC}"
echo -e "${BLUE}════════════════════════════════════════════════${NC}"
echo ""

# Función para check exitoso
check_pass() {
    echo -e "${GREEN}✅ $1${NC}"
    ((PASSED++))
}

# Función para check fallido
check_fail() {
    echo -e "${RED}❌ $1${NC}"
    ((FAILED++))
}

# Función para warning
check_warn() {
    echo -e "${YELLOW}⚠️  $1${NC}"
    ((WARNINGS++))
}

# 1. Verificar archivos esenciales
echo -e "${BLUE}📁 Verificando archivos...${NC}"

if [ -f "Dockerfile" ]; then
    check_pass "Dockerfile existe"
else
    check_fail "Dockerfile no encontrado"
fi

if [ -f "docker-compose.yml" ]; then
    check_pass "docker-compose.yml existe"
else
    check_warn "docker-compose.yml no encontrado (opcional)"
fi

if [ -f ".github/workflows/deploy.yml" ]; then
    check_pass "GitHub Actions workflow configurado"
else
    check_fail "Workflow de GitHub Actions no encontrado"
fi

if [ -f "go.mod" ]; then
    check_pass "go.mod existe"
else
    check_fail "go.mod no encontrado"
fi

if [ -d "migrations" ]; then
    migration_count=$(ls -1 migrations/*.sql 2>/dev/null | wc -l)
    if [ $migration_count -gt 0 ]; then
        check_pass "Migraciones encontradas ($migration_count archivos)"
    else
        check_warn "No hay archivos de migración SQL"
    fi
else
    check_fail "Directorio migrations/ no encontrado"
fi

echo ""

# 2. Verificar compilación
echo -e "${BLUE}🔨 Verificando compilación...${NC}"

if command -v go &> /dev/null; then
    check_pass "Go está instalado ($(go version | awk '{print $3}'))"
    
    echo -n "   Compilando... "
    if go build -o /tmp/fashion-blue-test ./cmd/api > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
        check_pass "Código compila sin errores"
        rm -f /tmp/fashion-blue-test
    else
        echo -e "${RED}✗${NC}"
        check_fail "Código no compila. Ejecuta: go build ./cmd/api"
    fi
else
    check_warn "Go no está instalado (no se puede verificar compilación)"
fi

echo ""

# 3. Verificar Docker
echo -e "${BLUE}🐳 Verificando Docker...${NC}"

if command -v docker &> /dev/null; then
    check_pass "Docker está instalado ($(docker --version | awk '{print $3}' | tr -d ','))"
    
    if docker info > /dev/null 2>&1; then
        check_pass "Docker daemon está corriendo"
    else
        check_warn "Docker daemon no está corriendo"
    fi
else
    check_warn "Docker no está instalado (necesario para deployment)"
fi

echo ""

# 4. Verificar Git
echo -e "${BLUE}🔧 Verificando Git...${NC}"

if command -v git &> /dev/null; then
    check_pass "Git está instalado"
    
    if git rev-parse --git-dir > /dev/null 2>&1; then
        check_pass "Repositorio Git inicializado"
        
        # Verificar remote
        if git remote get-url origin > /dev/null 2>&1; then
            remote_url=$(git remote get-url origin)
            check_pass "Remote configurado: $remote_url"
            
            # Verificar si es GitHub
            if [[ $remote_url == *"github.com"* ]]; then
                check_pass "Remote es GitHub (compatible con Actions)"
            else
                check_warn "Remote no es GitHub (GitHub Actions requiere GitHub)"
            fi
        else
            check_fail "No hay remote 'origin' configurado"
        fi
        
        # Verificar branch
        current_branch=$(git branch --show-current)
        if [ "$current_branch" == "main" ]; then
            check_pass "Branch actual es 'main'"
        else
            check_warn "Branch actual es '$current_branch' (workflow se ejecuta en 'main')"
        fi
        
        # Verificar cambios sin commit
        if git diff-index --quiet HEAD --; then
            check_pass "No hay cambios sin commit"
        else
            check_warn "Hay cambios sin commit"
        fi
    else
        check_fail "No estás en un repositorio Git"
    fi
else
    check_fail "Git no está instalado"
fi

echo ""

# 5. Verificar variables de entorno de ejemplo
echo -e "${BLUE}🔑 Verificando configuración...${NC}"

if [ -f ".env.example" ]; then
    check_pass ".env.example existe"
    
    # Verificar variables críticas en .env.example
    required_vars=("JWT_SECRET" "DB_HOST" "DB_USER" "DB_PASSWORD" "DB_NAME")
    for var in "${required_vars[@]}"; do
        if grep -q "^${var}=" .env.example 2>/dev/null; then
            check_pass "Variable $var presente en .env.example"
        else
            check_warn "Variable $var no encontrada en .env.example"
        fi
    done
else
    check_warn ".env.example no encontrado"
fi

echo ""

# 6. Verificar archivos de deployment específicos
echo -e "${BLUE}🚀 Verificando archivos de deployment...${NC}"

if [ -f "railway.json" ]; then
    check_pass "railway.json configurado"
else
    check_warn "railway.json no encontrado (solo necesario para Railway)"
fi

if [ -f ".do/app.yaml" ]; then
    check_pass ".do/app.yaml configurado"
else
    check_warn ".do/app.yaml no encontrado (solo necesario para DigitalOcean)"
fi

echo ""

# 7. Verificar .dockerignore
echo -e "${BLUE}📦 Verificando optimizaciones...${NC}"

if [ -f ".dockerignore" ]; then
    check_pass ".dockerignore existe"
else
    check_warn ".dockerignore no encontrado (recomendado para builds más rápidos)"
fi

# 8. Verificar .gitignore
if [ -f ".gitignore" ]; then
    check_pass ".gitignore existe"
    
    # Verificar que .env esté ignorado
    if grep -q "^\.env$" .gitignore 2>/dev/null; then
        check_pass ".env está en .gitignore (seguridad)"
    else
        check_fail ".env NO está en .gitignore (riesgo de seguridad)"
    fi
else
    check_warn ".gitignore no encontrado"
fi

echo ""

# 9. Verificar tamaño del proyecto
echo -e "${BLUE}📊 Estadísticas del proyecto...${NC}"

if [ -d ".git" ]; then
    repo_size=$(du -sh .git 2>/dev/null | awk '{print $1}')
    echo -e "   Tamaño del repositorio: ${repo_size}"
fi

go_files=$(find . -name "*.go" -not -path "./vendor/*" 2>/dev/null | wc -l)
echo -e "   Archivos Go: ${go_files}"

migration_files=$(find migrations -name "*.sql" 2>/dev/null | wc -l)
echo -e "   Archivos de migración: ${migration_files}"

echo ""

# Resumen final
echo -e "${BLUE}════════════════════════════════════════════════${NC}"
echo -e "${BLUE}                   RESUMEN                      ${NC}"
echo -e "${BLUE}════════════════════════════════════════════════${NC}"
echo ""
echo -e "  ${GREEN}✅ Checks pasados: ${PASSED}${NC}"
echo -e "  ${YELLOW}⚠️  Warnings: ${WARNINGS}${NC}"
echo -e "  ${RED}❌ Checks fallidos: ${FAILED}${NC}"
echo ""

# Recomendaciones
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}❌ No estás listo para deployment${NC}"
    echo -e "${YELLOW}Recomendaciones:${NC}"
    echo "   1. Corrige los errores mostrados arriba"
    echo "   2. Ejecuta este script nuevamente"
    echo "   3. Lee GITHUB_ACTIONS_DEPLOYMENT.md para más info"
    echo ""
    exit 1
elif [ $WARNINGS -gt 0 ]; then
    echo -e "${YELLOW}⚠️  Estás casi listo para deployment${NC}"
    echo -e "${YELLOW}Recomendaciones:${NC}"
    echo "   1. Revisa los warnings (no son críticos pero recomendados)"
    echo "   2. Lee QUICK_START_DEPLOYMENT.md para comenzar"
    echo ""
    exit 0
else
    echo -e "${GREEN}🎉 ¡Estás listo para deployment!${NC}"
    echo ""
    echo -e "${BLUE}Próximos pasos:${NC}"
    echo "   1. Lee: QUICK_START_DEPLOYMENT.md (deploy en 10 min)"
    echo "   2. O lee: GITHUB_ACTIONS_DEPLOYMENT.md (guía completa)"
    echo "   3. Configura GitHub Secrets para tu plataforma elegida"
    echo "   4. Push a main para deploy automático"
    echo ""
    exit 0
fi
