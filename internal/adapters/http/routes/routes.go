package routes

import (
	authHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/auth"
	capitalInjectionHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/capital_injection"
	categoryHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/category"
	customerHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/customer"
	paymentMethodHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/payment_method"
	productHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/product"
	sizeHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/size"
	supplierHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/supplier"
	userHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/user"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/middleware"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/auth"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/labstack/echo/v4"
)

// Handlers contiene todos los handlers HTTP
type Handlers struct {
	Auth             *authHandler.AuthHandler
	User             *userHandler.UserHandler
	Product          *productHandler.ProductHandler
	Category         *categoryHandler.CategoryHandler
	Size             *sizeHandler.SizeHandler
	PaymentMethod    *paymentMethodHandler.PaymentMethodHandler
	Customer         *customerHandler.CustomerHandler
	Supplier         *supplierHandler.SupplierHandler
	CapitalInjection *capitalInjectionHandler.CapitalInjectionHandler
}

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(e *echo.Echo, handlers Handlers, validateTokenUC *auth.ValidateTokenUseCase) {
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

	// Rutas protegidas - Inyecciones de Capital
	capitalInjections := api.Group("/capital-injections", authMiddleware)
	{
		capitalInjections.POST("", handlers.CapitalInjection.Create, middleware.RequireRole(entities.RoleSuperAdmin))
		capitalInjections.GET("", handlers.CapitalInjection.List)
		capitalInjections.GET("/:id", handlers.CapitalInjection.GetByID)
		capitalInjections.GET("/total", handlers.CapitalInjection.GetTotal)
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
}
