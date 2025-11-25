package swagger

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type SwaggerHandler struct {
	swaggerPath string
}

func NewSwaggerHandler(swaggerPath string) *SwaggerHandler {
	return &SwaggerHandler{
		swaggerPath: swaggerPath,
	}
}

// findSwaggerFile busca el archivo swagger.json en múltiples ubicaciones
func (h *SwaggerHandler) findSwaggerFile() ([]byte, error) {
	// Lista de rutas posibles
	possiblePaths := []string{
		h.swaggerPath,             // Ruta proporcionada
		"docs/swagger.json",       // Relativa al directorio actual
		"./docs/swagger.json",     // Explícitamente relativa
		"../docs/swagger.json",    // Un nivel arriba
		"../../docs/swagger.json", // Dos niveles arriba
	}

	// Intentar también con la ruta absoluta desde el working directory
	if wd, err := os.Getwd(); err == nil {
		possiblePaths = append(possiblePaths, filepath.Join(wd, "docs/swagger.json"))
	}

	// Intentar leer desde cada ruta
	for _, path := range possiblePaths {
		if data, err := os.ReadFile(path); err == nil {
			return data, nil
		}
	}

	return nil, os.ErrNotExist
}

// ServeSwaggerJSON sirve el archivo swagger.json
func (h *SwaggerHandler) ServeSwaggerJSON(c echo.Context) error {
	data, err := h.findSwaggerFile()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Swagger file not found. Run 'make swagger' to generate it.",
			"hint":  "Make sure docs/swagger.json exists in your project root",
		})
	}

	var swaggerSpec map[string]interface{}
	if err := json.Unmarshal(data, &swaggerSpec); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Invalid swagger file",
		})
	}

	return c.JSON(http.StatusOK, swaggerSpec)
}

// ServeSwaggerUI sirve la interfaz HTML de Swagger UI
func (h *SwaggerHandler) ServeSwaggerUI(c echo.Context) error {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Fashion Blue API - Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui.css">
    <style>
        body {
            margin: 0;
            padding: 0;
        }
        .swagger-ui .topbar {
            background-color: #1a1a1a;
        }
        .swagger-ui .topbar .download-url-wrapper {
            display: none;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "/api/v1/docs/swagger.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                defaultModelsExpandDepth: 1,
                defaultModelExpandDepth: 1,
                docExpansion: "list",
                filter: true,
                showExtensions: true,
                showCommonExtensions: true,
                tryItOutEnabled: true
            });
        };
    </script>
</body>
</html>
`
	return c.HTML(http.StatusOK, html)
}

// ServeRedocUI sirve la interfaz HTML de ReDoc (alternativa a Swagger UI)
func (h *SwaggerHandler) ServeRedocUI(c echo.Context) error {
	html := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Fashion Blue API - ReDoc</title>
    <style>
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
    <redoc spec-url='/api/v1/docs/swagger.json'></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
</body>
</html>
`
	return c.HTML(http.StatusOK, html)
}

// ServeRapidocUI sirve la interfaz HTML de RapiDoc (otra alternativa)
func (h *SwaggerHandler) ServeRapidocUI(c echo.Context) error {
	html := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Fashion Blue API - RapiDoc</title>
    <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
</head>
<body>
    <rapi-doc
        spec-url="/api/v1/docs/swagger.json"
        theme="dark"
        bg-color="#1a1a1a"
        text-color="#ffffff"
        primary-color="#4CAF50"
        render-style="read"
        show-header="true"
        show-info="true"
        allow-authentication="true"
        allow-server-selection="true"
        allow-api-list-style-selection="true"
        default-api-server="http://localhost:8080"
    >
    </rapi-doc>
</body>
</html>
`
	return c.HTML(http.StatusOK, html)
}
