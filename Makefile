.PHONY: help run build test clean docker-up docker-down migrate-up migrate-down migrate-create

help: ## Mostrar ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run: ## Ejecutar la aplicación
	go run cmd/api/main.go

build: ## Compilar la aplicación
	go build -o bin/fashion-blue cmd/api/main.go

test: ## Ejecutar tests
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-coverage: test ## Ejecutar tests y mostrar cobertura
	go tool cover -html=coverage.txt -o coverage.html
	open coverage.html

clean: ## Limpiar archivos generados
	rm -rf bin/
	rm -f coverage.txt coverage.html

docker-up: ## Levantar contenedores Docker
	docker-compose up -d

docker-down: ## Detener contenedores Docker
	docker-compose down

docker-logs: ## Ver logs de Docker
	docker-compose logs -f

docker-rebuild: ## Reconstruir y levantar contenedores
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

migrate-up: ## Ejecutar migraciones pendientes
	migrate -path migrations -database "postgresql://fashionblue:fashionblue123@localhost:5432/fashionblue_db?sslmode=disable" up

migrate-down: ## Revertir última migración
	migrate -path migrations -database "postgresql://fashionblue:fashionblue123@localhost:5432/fashionblue_db?sslmode=disable" down 1

migrate-create: ## Crear nueva migración (uso: make migrate-create name=nombre_migracion)
	migrate create -ext sql -dir migrations -seq $(name)

deps: ## Instalar dependencias
	go mod download
	go mod tidy

lint: ## Ejecutar linter
	golangci-lint run

format: ## Formatear código
	go fmt ./...
	goimports -w .

swagger: ## Generar documentación Swagger
	@echo "🔍 Generando documentación Swagger..."
	@go run scripts/generate_swagger.go
	@echo "✅ Documentación generada en docs/swagger.json"

swagger-ui: swagger ## Generar Swagger y abrir en navegador
	@echo "🌐 Abriendo Swagger UI..."
	@open https://editor.swagger.io/

.DEFAULT_GOAL := help
