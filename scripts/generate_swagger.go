package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// SwaggerSpec representa la especificación OpenAPI 3.0
type SwaggerSpec struct {
	OpenAPI    string                 `json:"openapi"`
	Info       SwaggerInfo            `json:"info"`
	Servers    []SwaggerServer        `json:"servers"`
	Paths      map[string]interface{} `json:"paths"`
	Components SwaggerComponents      `json:"components"`
}

type SwaggerInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type SwaggerServer struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type SwaggerComponents struct {
	Schemas         map[string]interface{} `json:"schemas"`
	SecuritySchemes map[string]interface{} `json:"securitySchemes"`
}

type EndpointInfo struct {
	Method      string
	Path        string
	Handler     string
	Summary     string
	Description string
	RequestBody string
	Response    string
}

func main3() {
	fmt.Println("🔍 Generando documentación Swagger...")

	// Crear spec base
	spec := SwaggerSpec{
		OpenAPI: "3.0.0",
		Info: SwaggerInfo{
			Title:       "Fashion Blue API",
			Description: "API para gestión de órdenes, productos, clientes y más",
			Version:     "1.0.0",
		},
		Servers: []SwaggerServer{
			{
				URL:         "http://localhost:8080",
				Description: "Servidor de desarrollo",
			},
		},
		Paths: make(map[string]interface{}),
		Components: SwaggerComponents{
			Schemas: make(map[string]interface{}),
			SecuritySchemes: map[string]interface{}{
				"bearerAuth": map[string]interface{}{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
				},
			},
		},
	}

	// Analizar rutas
	endpoints := analyzeRoutes("internal/adapters/http/routes/routes.go")

	// Analizar handlers
	analyzeHandlers("internal/adapters/http/handlers", &spec, endpoints)

	// Guardar spec
	saveSwaggerSpec(&spec, "docs/swagger.json")

	fmt.Println("✅ Documentación Swagger generada en docs/swagger.json")
	fmt.Println("📝 Para visualizar: https://editor.swagger.io/")
}

func analyzeRoutes(routesFile string) []EndpointInfo {
	endpoints := []EndpointInfo{}
	groupPaths := make(map[string]string) // Mapeo de variable de grupo a su path base

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, routesFile, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("⚠️  Error parsing routes: %v\n", err)
		return endpoints
	}

	// Primero, encontrar todas las declaraciones de grupos
	ast.Inspect(node, func(n ast.Node) bool {
		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		// Buscar asignaciones como: customers := api.Group("/customers", ...)
		if len(assignStmt.Lhs) == 1 && len(assignStmt.Rhs) == 1 {
			ident, ok := assignStmt.Lhs[0].(*ast.Ident)
			if !ok {
				return true
			}

			callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
			if !ok {
				return true
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok || selExpr.Sel.Name != "Group" {
				return true
			}

			// Extraer el path del grupo
			if len(callExpr.Args) >= 1 {
				groupPath := extractStringLiteral(callExpr.Args[0])
				if groupPath != "" {
					groupPaths[ident.Name] = groupPath
				}
			}
		}

		return true
	})

	// Luego, buscar llamadas a métodos HTTP
	ast.Inspect(node, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Verificar si es una llamada a método HTTP
		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		method := selExpr.Sel.Name
		if !isHTTPMethod(method) {
			return true
		}

		// Extraer el nombre del grupo (ej: "customers" en customers.GET(...))
		groupIdent, ok := selExpr.X.(*ast.Ident)
		if !ok {
			return true
		}

		groupBasePath := groupPaths[groupIdent.Name]

		// Extraer path y handler
		if len(callExpr.Args) >= 2 {
			routePath := extractStringLiteral(callExpr.Args[0])
			handler := extractHandlerName(callExpr.Args[1])

			if handler != "" {
				// Combinar el path base del grupo con el path de la ruta
				fullPath := groupBasePath + routePath

				// Excluir rutas de documentación
				if strings.Contains(fullPath, "/docs") ||
					strings.Contains(handler, "Swagger") ||
					strings.Contains(handler, "Redoc") ||
					strings.Contains(handler, "Rapidoc") {
					return true
				}

				endpoints = append(endpoints, EndpointInfo{
					Method:  strings.ToUpper(method),
					Path:    fullPath,
					Handler: handler,
				})
			}
		}

		return true
	})

	return endpoints
}

func analyzeHandlers(handlersDir string, spec *SwaggerSpec, endpoints []EndpointInfo) {
	filepath.Walk(handlersDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, "_handler.go") {
			return nil
		}

		analyzeHandlerFile(path, spec, endpoints)
		return nil
	})
}

func analyzeHandlerFile(filePath string, spec *SwaggerSpec, endpoints []EndpointInfo) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return
	}

	// Analizar cada función del handler
	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil {
			continue
		}

		handlerName := funcDecl.Name.Name

		// Buscar endpoint correspondiente
		for _, endpoint := range endpoints {
			if strings.Contains(endpoint.Handler, handlerName) {
				addEndpointToSpec(spec, endpoint, funcDecl)
				break
			}
		}
	}
}

func addEndpointToSpec(spec *SwaggerSpec, endpoint EndpointInfo, funcDecl *ast.FuncDecl) {
	// Extraer tag del path ORIGINAL (antes de convertir)
	tag := extractTag(endpoint.Path)

	// Convertir path de Echo (:id) a OpenAPI ({id})
	path := convertPathToOpenAPI(endpoint.Path)

	// Agregar prefijo /api/v1 si no lo tiene
	if !strings.HasPrefix(path, "/api/v1") {
		path = "/api/v1" + path
	}

	method := strings.ToLower(endpoint.Method)

	// Crear path si no existe
	if spec.Paths[path] == nil {
		spec.Paths[path] = make(map[string]interface{})
	}

	pathItem := spec.Paths[path].(map[string]interface{})

	// Extraer documentación del comentario
	summary := funcDecl.Name.Name
	description := ""
	requestDTO := ""
	responseDTO := ""

	if funcDecl.Doc != nil {
		for _, comment := range funcDecl.Doc.List {
			text := strings.TrimPrefix(comment.Text, "//")
			text = strings.TrimSpace(text)

			// Buscar anotaciones especiales
			if strings.HasPrefix(text, "@Request:") {
				requestDTO = strings.TrimSpace(strings.TrimPrefix(text, "@Request:"))
			} else if strings.HasPrefix(text, "@Response:") {
				responseDTO = strings.TrimSpace(strings.TrimPrefix(text, "@Response:"))
			} else if description == "" && !strings.HasPrefix(text, "@") {
				description = text
			}
		}
	}

	// Si no hay @Request, intentar extraer del Bind()
	if requestDTO == "" {
		requestDTO = extractBindDTO(funcDecl)
	}

	// Crear operación
	operation := map[string]interface{}{
		"summary":     summary,
		"description": description,
		"tags":        []string{tag},
		"responses":   createResponses(responseDTO),
	}

	// Agregar parámetros de path (buscar tanto : como {})
	if strings.Contains(path, "{") || strings.Contains(path, ":") {
		params := extractPathParams(path)
		if len(params) > 0 {
			operation["parameters"] = params
		}
	}

	// Agregar request body para POST/PUT/PATCH
	if method == "post" || method == "put" || method == "patch" {
		schema := createSchemaForDTO(requestDTO, path)

		operation["requestBody"] = map[string]interface{}{
			"required": true,
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": schema,
				},
			},
		}
	}

	// Agregar seguridad (JWT) solo si NO es una ruta pública
	isPublicRoute := strings.Contains(path, "/auth/login") ||
		strings.Contains(path, "/auth/register")

	if !isPublicRoute {
		operation["security"] = []map[string][]string{
			{"bearerAuth": {}},
		}
	}

	pathItem[method] = operation
}

// createResponses crea las respuestas basadas en el DTO de respuesta
func createResponses(responseDTO string) map[string]interface{} {
	responses := map[string]interface{}{
		"400": map[string]interface{}{
			"description": "Bad request",
		},
		"401": map[string]interface{}{
			"description": "Unauthorized",
		},
		"500": map[string]interface{}{
			"description": "Internal server error",
		},
	}

	// Response 200 con schema específico si hay DTO
	if responseDTO != "" {
		schema := createSchemaForDTO(responseDTO, "")
		responses["200"] = map[string]interface{}{
			"description": "Successful response",
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": schema,
				},
			},
		}
	} else {
		// Response genérico
		responses["200"] = map[string]interface{}{
			"description": "Successful response",
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"success": map[string]interface{}{"type": "boolean"},
							"message": map[string]interface{}{"type": "string"},
							"data":    map[string]interface{}{"type": "object"},
						},
					},
				},
			},
		}
	}

	return responses
}

// createSchemaForDTO crea un schema basado en el nombre del DTO detectado
func createSchemaForDTO(dtoName string, path string) map[string]interface{} {
	schema := map[string]interface{}{
		"type": "object",
	}

	// Mapeo de DTOs conocidos a sus propiedades
	dtoSchemas := map[string]map[string]interface{}{
		"LoginRequest": {
			"properties": map[string]interface{}{
				"email":    map[string]interface{}{"type": "string", "example": "admin@example.com"},
				"password": map[string]interface{}{"type": "string", "example": "password123"},
			},
			"required": []string{"email", "password"},
		},
		"RegisterRequest": {
			"properties": map[string]interface{}{
				"name":     map[string]interface{}{"type": "string", "example": "John Doe"},
				"email":    map[string]interface{}{"type": "string", "example": "user@example.com"},
				"password": map[string]interface{}{"type": "string", "example": "password123"},
				"role":     map[string]interface{}{"type": "string", "example": "admin"},
			},
			"required": []string{"name", "email", "password", "role"},
		},
		"Product": {
			"properties": map[string]interface{}{
				"name":        map[string]interface{}{"type": "string", "example": "Camiseta Básica"},
				"description": map[string]interface{}{"type": "string", "example": "Camiseta de algodón"},
				"category_id": map[string]interface{}{"type": "integer", "example": 1},
				"price":       map[string]interface{}{"type": "number", "example": 25000},
				"cost":        map[string]interface{}{"type": "number", "example": 15000},
				"stock":       map[string]interface{}{"type": "integer", "example": 50},
			},
			"required": []string{"name", "category_id", "price", "cost"},
		},
		"Category": {
			"properties": map[string]interface{}{
				"name":        map[string]interface{}{"type": "string", "example": "Camisetas"},
				"description": map[string]interface{}{"type": "string", "example": "Categoría de camisetas"},
			},
			"required": []string{"name"},
		},
		"Customer": {
			"properties": map[string]interface{}{
				"name":    map[string]interface{}{"type": "string", "example": "Juan Pérez"},
				"email":   map[string]interface{}{"type": "string", "example": "juan@example.com"},
				"phone":   map[string]interface{}{"type": "string", "example": "3001234567"},
				"address": map[string]interface{}{"type": "string", "example": "Calle 123 #45-67"},
			},
			"required": []string{"name"},
		},
		"Supplier": {
			"properties": map[string]interface{}{
				"name":    map[string]interface{}{"type": "string", "example": "Proveedor XYZ"},
				"contact": map[string]interface{}{"type": "string", "example": "Carlos López"},
				"phone":   map[string]interface{}{"type": "string", "example": "3009876543"},
				"email":   map[string]interface{}{"type": "string", "example": "proveedor@example.com"},
			},
			"required": []string{"name"},
		},
		"Order": {
			"properties": map[string]interface{}{
				"customer_id": map[string]interface{}{"type": "integer", "example": 1},
				"type":        map[string]interface{}{"type": "string", "example": "READY_MADE"},
				"notes":       map[string]interface{}{"type": "string", "example": "Entrega urgente"},
			},
			"required": []string{"customer_id", "type"},
		},
		"OrderItem": {
			"properties": map[string]interface{}{
				"product_id": map[string]interface{}{"type": "integer", "example": 1},
				"quantity":   map[string]interface{}{"type": "integer", "example": 2},
				"size_id":    map[string]interface{}{"type": "integer", "example": 1},
			},
			"required": []string{"product_id", "quantity"},
		},
		"Payment": {
			"properties": map[string]interface{}{
				"amount":            map[string]interface{}{"type": "number", "example": 50000},
				"payment_method_id": map[string]interface{}{"type": "integer", "example": 1},
				"notes":             map[string]interface{}{"type": "string", "example": "Abono parcial"},
			},
			"required": []string{"amount", "payment_method_id"},
		},
		"Transaction": {
			"properties": map[string]interface{}{
				"customer_id": map[string]interface{}{"type": "integer", "example": 1},
				"amount":      map[string]interface{}{"type": "number", "example": 100000},
				"type":        map[string]interface{}{"type": "string", "example": "CREDIT"},
				"description": map[string]interface{}{"type": "string", "example": "Ajuste manual"},
			},
			"required": []string{"customer_id", "amount", "type"},
		},
		"CapitalInjection": {
			"properties": map[string]interface{}{
				"amount":      map[string]interface{}{"type": "number", "example": 1000000},
				"description": map[string]interface{}{"type": "string", "example": "Inversión inicial"},
				"date":        map[string]interface{}{"type": "string", "example": "2024-01-15"},
			},
			"required": []string{"amount", "description"},
		},
		"ChangePasswordRequest": {
			"properties": map[string]interface{}{
				"current_password": map[string]interface{}{"type": "string", "example": "oldpass123"},
				"new_password":     map[string]interface{}{"type": "string", "example": "newpass123"},
			},
			"required": []string{"current_password", "new_password"},
		},
		"UpdateOrderStatusRequest": {
			"properties": map[string]interface{}{
				"status": map[string]interface{}{"type": "string", "example": "APPROVED"},
			},
			"required": []string{"status"},
		},
	}

	// Si encontramos el DTO en nuestro mapeo, usar su schema
	if dtoSchema, ok := dtoSchemas[dtoName]; ok {
		for k, v := range dtoSchema {
			schema[k] = v
		}
	} else if dtoName != "" {
		// Si tenemos un DTO pero no está mapeado, al menos indicar su nombre
		schema["description"] = "Request body for " + dtoName
		schema["example"] = map[string]interface{}{
			"field": "value",
		}
	} else {
		// Sin DTO detectado, schema genérico
		schema["example"] = map[string]interface{}{
			"field": "value",
		}
	}

	return schema
}

// extractBindDTO extrae el nombre del DTO usado en c.Bind(&dto)
func extractBindDTO(funcDecl *ast.FuncDecl) string {
	if funcDecl.Body == nil {
		return ""
	}

	var dtoName string

	// Primero buscar declaraciones var con tipo explícito (var req LoginRequest)
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		declStmt, ok := n.(*ast.DeclStmt)
		if !ok {
			return true
		}

		genDecl, ok := declStmt.Decl.(*ast.GenDecl)
		if !ok {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			// Obtener el tipo de la variable
			if valueSpec.Type != nil {
				if ident, ok := valueSpec.Type.(*ast.Ident); ok {
					// Tipo simple como LoginRequest
					dtoName = ident.Name
					return false
				} else if selExpr, ok := valueSpec.Type.(*ast.SelectorExpr); ok {
					// Tipo calificado como dto.Product
					dtoName = selExpr.Sel.Name
					return false
				}
			}
		}

		return true
	})

	// Si no encontramos nada, buscar en asignaciones con :=
	if dtoName == "" {
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			assignStmt, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}

			if len(assignStmt.Lhs) > 0 && len(assignStmt.Rhs) > 0 {
				if compLit, ok := assignStmt.Rhs[0].(*ast.CompositeLit); ok {
					if selExpr, ok := compLit.Type.(*ast.SelectorExpr); ok {
						dtoName = selExpr.Sel.Name
						return false
					} else if ident, ok := compLit.Type.(*ast.Ident); ok {
						dtoName = ident.Name
						return false
					}
				}
			}

			return true
		})
	}

	return dtoName
}

func extractTag(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	// Buscar la primera parte que NO sea un parámetro (no empieza con :)
	for _, part := range parts {
		// Saltar partes vacías, parámetros (:id), y prefijos comunes
		if part == "" || strings.HasPrefix(part, ":") || part == "api" || part == "v1" {
			continue
		}
		// Retornar la primera parte válida con primera letra mayúscula
		return strings.Title(part)
	}

	return "General"
}

func convertPathToOpenAPI(path string) string {
	// Convertir :param a {param}
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			paramName := strings.TrimPrefix(part, ":")
			parts[i] = "{" + paramName + "}"
		}
	}
	return strings.Join(parts, "/")
}

func extractPathParams(path string) []map[string]interface{} {
	params := []map[string]interface{}{}
	parts := strings.Split(path, "/")

	for _, part := range parts {
		// Detectar tanto :param como {param}
		paramName := ""
		if strings.HasPrefix(part, ":") {
			paramName = strings.TrimPrefix(part, ":")
		} else if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			paramName = strings.Trim(part, "{}")
		}

		if paramName != "" {
			params = append(params, map[string]interface{}{
				"name":        paramName,
				"in":          "path",
				"required":    true,
				"schema":      map[string]string{"type": "string"},
				"description": fmt.Sprintf("ID of %s", paramName),
			})
		}
	}

	return params
}

func isHTTPMethod(method string) bool {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	upper := strings.ToUpper(method)
	for _, m := range methods {
		if m == upper {
			return true
		}
	}
	return false
}

func extractStringLiteral(expr ast.Expr) string {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		return strings.Trim(lit.Value, `"`)
	}
	return ""
}

func extractHandlerName(expr ast.Expr) string {
	if selExpr, ok := expr.(*ast.SelectorExpr); ok {
		return selExpr.Sel.Name
	}
	return ""
}

func saveSwaggerSpec(spec *SwaggerSpec, outputPath string) {
	// Crear directorio si no existe
	dir := filepath.Dir(outputPath)
	os.MkdirAll(dir, 0755)

	// Serializar a JSON
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		fmt.Printf("❌ Error serializing spec: %v\n", err)
		return
	}

	// Guardar archivo
	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing file: %v\n", err)
		return
	}
}
