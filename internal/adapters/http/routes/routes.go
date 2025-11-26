package routes

import (
	analyticsHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/analytics"
	auditHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/audit"
	authHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/auth"
	categoryHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/category"
	customerHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/customer"
	financialTransactionHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/financial_transaction"
	orderHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/order"
	paymentMethodHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/payment_method"
	productHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/product"
	sizeHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/size"
	supplierHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/supplier"
	swaggerHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/swagger"
	userHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/user"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/middleware"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/auth"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/labstack/echo/v4"
)

// Handlers contiene todos los handlers HTTP
type Handlers struct {
	Analytics            *analyticsHandler.AnalyticsHTTPHandler
	Audit                *auditHandler.AuditHTTPHandler
	Auth                 *authHandler.AuthHandler
	User                 *userHandler.UserHandler
	Product              *productHandler.ProductHandler
	Category             *categoryHandler.CategoryHandler
	Size                 *sizeHandler.SizeHandler
	PaymentMethod        *paymentMethodHandler.PaymentMethodHandler
	Customer             *customerHandler.CustomerHandler
	CustomerStatement    *customerHandler.StatementHandler
	Order                *orderHandler.OrderHandler
	Supplier             *supplierHandler.SupplierHandler
	FinancialTransaction *financialTransactionHandler.FinancialTransactionHandler
	Swagger              *swaggerHandler.SwaggerHandler
}

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(e *echo.Echo, handlers Handlers, validateTokenUC *auth.ValidateTokenUseCase) {
	// Health check endpoint (para Railway, Docker, K8s, etc.)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":  "ok",
			"service": "fashion-blue-api",
			"version": "1.0.0",
		})
	})

	// API v1
	api := e.Group("/api/v1")

	// Rutas públicas - Autenticación
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", handlers.Auth.Login)
		authGroup.POST("/register", handlers.Auth.Register)
	}

	// Middleware de autenticación para rutas protegidas
	authMiddleware := middleware.AuthMiddleware(validateTokenUC)

	// Rutas protegidas - Categorías
	categories := api.Group("/categories", authMiddleware)
	{
		categories.POST("", handlers.Category.Create, middleware.RequireRole(entities.RoleSuperAdmin))
		categories.GET("", handlers.Category.List)
		categories.GET("/:id", handlers.Category.GetByID)
		categories.PUT("/:id", handlers.Category.Update, middleware.RequireRole(entities.RoleSuperAdmin))
		categories.DELETE("/:id", handlers.Category.Delete, middleware.RequireRole(entities.RoleSuperAdmin))
	}

	// Rutas protegidas - Productos
	products := api.Group("/products", authMiddleware)
	{
		products.POST("", handlers.Product.Create, middleware.RequireRole(entities.RoleSuperAdmin))
		products.GET("", handlers.Product.List)
		products.GET("/:id", handlers.Product.GetByID)
		products.GET("/low-stock", handlers.Product.GetLowStock)
		products.PUT("/:id", handlers.Product.Update, middleware.RequireRole(entities.RoleSuperAdmin))
		products.DELETE("/:id", handlers.Product.Delete, middleware.RequireRole(entities.RoleSuperAdmin))

		// Rutas de fotos de productos
		products.POST("/:id/photos", handlers.Product.UploadPhotos)
		products.GET("/:id/photos", handlers.Product.GetPhotos)
		products.DELETE("/:id/photos/:photoId", handlers.Product.DeletePhoto)
		products.PUT("/:id/photos/:photoId/primary", handlers.Product.SetPrimaryPhoto)
	}

	// Rutas protegidas - Tallas
	sizes := api.Group("/sizes", authMiddleware)
	{
		sizes.GET("", handlers.Size.List)
		sizes.GET("/:id", handlers.Size.GetByID)
		sizes.GET("/type/:type", handlers.Size.GetByType)
	}

	// Rutas protegidas - Métodos de Pago
	paymentMethods := api.Group("/payment-methods", authMiddleware)
	{
		paymentMethods.GET("", handlers.PaymentMethod.List)
	}

	// Rutas protegidas - Clientes
	customers := api.Group("/customers", authMiddleware)
	{
		customers.POST("", handlers.Customer.Create)
		customers.GET("", handlers.Customer.List)
		customers.POST("/transactions", handlers.Customer.AddTransaction)          // Nuevo endpoint para movimientos manuales
		customers.GET("/upcoming-payments", handlers.Customer.GetUpcomingPayments) // Debe ir antes de /:id
		customers.GET("/:id", handlers.Customer.GetByID)
		customers.GET("/:id/balance", handlers.Customer.GetBalance)
		customers.GET("/:id/history", handlers.Customer.GetHistory)
		customers.GET("/:id/statement", handlers.CustomerStatement.DownloadStatement) // PDF estado de cuenta (days opcional)
		customers.POST("/:id/payments", handlers.Customer.CreatePayment)
		customers.PUT("/:id", handlers.Customer.Update)
		customers.DELETE("/:id", handlers.Customer.Delete, middleware.RequireRole(entities.RoleSuperAdmin))
	}

	// Rutas protegidas - Proveedores
	suppliers := api.Group("/suppliers", authMiddleware)
	{
		suppliers.POST("", handlers.Supplier.Create, middleware.RequireRole(entities.RoleSuperAdmin))
		suppliers.GET("", handlers.Supplier.List)
		suppliers.GET("/:id", handlers.Supplier.GetByID)
		suppliers.PUT("/:id", handlers.Supplier.Update, middleware.RequireRole(entities.RoleSuperAdmin))
		suppliers.DELETE("/:id", handlers.Supplier.Delete, middleware.RequireRole(entities.RoleSuperAdmin))
	}

	// Rutas protegidas - Transacciones Financieras (Solo Super Admin)
	financialTransactions := api.Group("/financial-transactions", authMiddleware, middleware.RequireRole(entities.RoleSuperAdmin))
	{
		financialTransactions.POST("", handlers.FinancialTransaction.Create)            // Crear ingreso o gasto
		financialTransactions.GET("", handlers.FinancialTransaction.List)               // Listar con filtros
		financialTransactions.GET("/balance", handlers.FinancialTransaction.GetBalance) // Obtener balance
		financialTransactions.GET("/:id", handlers.FinancialTransaction.GetByID)        // Obtener por ID
	}

	// Rutas protegidas - Órdenes
	orders := api.Group("/orders", authMiddleware)
	{
		orders.POST("", handlers.Order.CreateOrder)
		orders.GET("", handlers.Order.ListOrders)
		orders.GET("/:id", handlers.Order.GetOrder)
		orders.GET("/:id/allowed-statuses", handlers.Order.GetAllowedNextStatuses) // Obtener estados permitidos
		orders.POST("/:id/change-status", handlers.Order.ChangeOrderStatus)        // Cambiar estado
		orders.POST("/:id/items", handlers.Order.AddOrderItem)
		orders.PUT("/:id/items/:itemId", handlers.Order.UpdateOrderItem)
		orders.DELETE("/:id/items/:itemId", handlers.Order.RemoveOrderItem)
	}

	// Rutas protegidas - Usuarios (Solo Super Admin)
	users := api.Group("/users", authMiddleware, middleware.RequireRole(entities.RoleSuperAdmin))
	{
		users.POST("", handlers.User.Create)
		users.GET("", handlers.User.List)
		users.GET("/:id", handlers.User.GetByID)
		users.PUT("/:id", handlers.User.Update)
		users.DELETE("/:id", handlers.User.Delete)
		users.PUT("/:id/password", handlers.User.ChangePassword)
	}

	// Rutas de Analytics (protegidas - solo admin)
	analytics := api.Group("/analytics")
	analytics.Use(authMiddleware)
	analytics.Use(middleware.RequireRole(entities.RoleSuperAdmin))
	{
		analytics.GET("/metrics", handlers.Analytics.GetMetrics)
		analytics.GET("/dashboard", handlers.Analytics.GetDashboardSummary)
	}

	// Rutas de Auditoría (protegidas - solo admin)
	audit := api.Group("/audit")
	audit.Use(authMiddleware)
	audit.Use(middleware.RequireRole(entities.RoleSuperAdmin))
	{
		audit.GET("/logs", handlers.Audit.GetAuditLogs)
		audit.GET("/logs/:orderId", handlers.Audit.GetAuditLogsByOrder) // Busca por Order ID
		audit.GET("/users/:userId/logs", handlers.Audit.GetAuditLogsByUser)
		audit.GET("/stats", handlers.Audit.GetAuditStats)
	}

	// Rutas públicas - Documentación API (Swagger)
	docs := api.Group("/docs")
	{
		docs.GET("", handlers.Swagger.ServeSwaggerUI)                // Swagger UI
		docs.GET("/swagger.json", handlers.Swagger.ServeSwaggerJSON) // JSON spec
		docs.GET("/redoc", handlers.Swagger.ServeRedocUI)            // ReDoc UI
		docs.GET("/rapidoc", handlers.Swagger.ServeRapidocUI)        // RapiDoc UI
	}
}
