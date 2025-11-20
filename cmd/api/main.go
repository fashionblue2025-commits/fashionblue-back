package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	authHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/auth"
	capitalInjectionHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/capital_injection"
	categoryHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/category"
	customerHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/customer"
	paymentMethodHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/payment_method"
	productHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/product"
	sizeHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/size"
	supplierHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/supplier"
	userHandler "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/user"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/routes"
	capitalInjectionRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/capital_injection"
	categoryRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/category"
	customerRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/customer"
	paymentMethodRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/payment_method"
	productRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/product"
	sizeRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/size"
	supplierRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/supplier"
	userRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/user"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/auth"
	capitalInjectionUseCases "github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/capital_injection"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/category"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/customer"
	paymentMethodUseCases "github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/payment_method"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/product"
	sizeUseCases "github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/size"
	supplierUseCases "github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/supplier"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/user"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/config"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Conectar a la base de datos
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Ejecutar migraciones automáticas
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Inicializar repositorios
	userRepository := userRepo.NewUserRepository(db)
	categoryRepository := categoryRepo.NewCategoryRepository(db)
	productRepository := productRepo.NewProductRepository(db)
	sizeRepository := sizeRepo.NewSizeRepository(db)
	paymentMethodRepository := paymentMethodRepo.NewPaymentMethodRepository(db)
	customerRepository := customerRepo.NewCustomerRepository(db)
	customerTransactionRepository := customerRepo.NewCustomerTransactionRepository(db)
	supplierRepository := supplierRepo.NewSupplierRepository(db)
	capitalInjectionRepository := capitalInjectionRepo.NewCapitalInjectionRepository(db)

	// Inicializar casos de uso - Auth
	loginUC := auth.NewLoginUseCase(userRepository, cfg.JWT.Secret, cfg.JWT.GetExpiration())
	registerUC := auth.NewRegisterUseCase(userRepository, cfg.JWT.Secret, cfg.JWT.GetExpiration())
	validateTokenUC := auth.NewValidateTokenUseCase(userRepository, cfg.JWT.Secret)

	// Inicializar casos de uso - User
	createUserUC := user.NewCreateUserUseCase(userRepository)
	getUserUC := user.NewGetUserUseCase(userRepository)
	listUsersUC := user.NewListUsersUseCase(userRepository)
	updateUserUC := user.NewUpdateUserUseCase(userRepository)
	deleteUserUC := user.NewDeleteUserUseCase(userRepository)
	changePasswordUC := user.NewChangePasswordUseCase(userRepository)

	// Inicializar casos de uso - Product
	createProductUC := product.NewCreateProductUseCase(productRepository)
	getProductUC := product.NewGetProductUseCase(productRepository)
	listProductsUC := product.NewListProductsUseCase(productRepository)
	updateProductUC := product.NewUpdateProductUseCase(productRepository)
	deleteProductUC := product.NewDeleteProductUseCase(productRepository)
	getLowStockUC := product.NewGetLowStockProductsUseCase(productRepository)

	// Inicializar casos de uso - Category
	createCategoryUC := category.NewCreateCategoryUseCase(categoryRepository)
	getCategoryUC := category.NewGetCategoryUseCase(categoryRepository)
	listCategoriesUC := category.NewListCategoriesUseCase(categoryRepository)
	updateCategoryUC := category.NewUpdateCategoryUseCase(categoryRepository)
	deleteCategoryUC := category.NewDeleteCategoryUseCase(categoryRepository)

	// Inicializar casos de uso - Size
	listSizesUC := sizeUseCases.NewListSizesUseCase(sizeRepository)
	getSizeUC := sizeUseCases.NewGetSizeUseCase(sizeRepository)
	getSizesByTypeUC := sizeUseCases.NewGetSizesByTypeUseCase(sizeRepository)

	// Inicializar casos de uso - PaymentMethod
	listPaymentMethodsUC := paymentMethodUseCases.NewListPaymentMethodsUseCase(paymentMethodRepository)

	// Inicializar casos de uso - Customer
	createCustomerUC := customer.NewCreateCustomerUseCase(customerRepository)
	getCustomerUC := customer.NewGetCustomerUseCase(customerRepository)
	listCustomersUC := customer.NewListCustomersUseCase(customerRepository)
	updateCustomerUC := customer.NewUpdateCustomerUseCase(customerRepository)
	deleteCustomerUC := customer.NewDeleteCustomerUseCase(customerRepository)
	getCustomerHistoryUC := customer.NewGetCustomerHistoryUseCase(customerTransactionRepository)
	createPaymentUC := customer.NewCreatePaymentUseCase(customerTransactionRepository, customerRepository)
	getUpcomingPaymentsUC := customer.NewGetUpcomingPaymentsUseCase(customerRepository)
	getCustomerBalanceUC := customer.NewGetCustomerBalanceUseCase(customerRepository)
	addTransactionUC := customer.NewAddTransactionUseCase(customerTransactionRepository, customerRepository)

	// Inicializar casos de uso - Supplier
	createSupplierUC := supplierUseCases.NewCreateSupplierUseCase(supplierRepository)
	getSupplierUC := supplierUseCases.NewGetSupplierUseCase(supplierRepository)
	listSuppliersUC := supplierUseCases.NewListSuppliersUseCase(supplierRepository)
	updateSupplierUC := supplierUseCases.NewUpdateSupplierUseCase(supplierRepository)
	deleteSupplierUC := supplierUseCases.NewDeleteSupplierUseCase(supplierRepository)

	// Inicializar casos de uso - CapitalInjection
	createInjectionUC := capitalInjectionUseCases.NewCreateInjectionUseCase(capitalInjectionRepository)
	getInjectionUC := capitalInjectionUseCases.NewGetInjectionUseCase(capitalInjectionRepository)
	listInjectionsUC := capitalInjectionUseCases.NewListInjectionsUseCase(capitalInjectionRepository)
	getTotalCapitalUC := capitalInjectionUseCases.NewGetTotalCapitalUseCase(capitalInjectionRepository)

	// Inicializar handlers
	authHandlerInstance := authHandler.NewAuthHandler(loginUC, registerUC)
	userHandlerInstance := userHandler.NewUserHandler(createUserUC, getUserUC, listUsersUC, updateUserUC, deleteUserUC, changePasswordUC)
	productHandlerInstance := productHandler.NewProductHandler(createProductUC, getProductUC, listProductsUC, updateProductUC, deleteProductUC, getLowStockUC)
	categoryHandlerInstance := categoryHandler.NewCategoryHandler(createCategoryUC, getCategoryUC, listCategoriesUC, updateCategoryUC, deleteCategoryUC)
	sizeHandlerInstance := sizeHandler.NewSizeHandler(listSizesUC, getSizeUC, getSizesByTypeUC)
	paymentMethodHandlerInstance := paymentMethodHandler.NewPaymentMethodHandler(listPaymentMethodsUC)
	customerHandlerInstance := customerHandler.NewCustomerHandler(createCustomerUC, getCustomerUC, listCustomersUC, updateCustomerUC, deleteCustomerUC, getCustomerHistoryUC, createPaymentUC, getUpcomingPaymentsUC, getCustomerBalanceUC, addTransactionUC)
	supplierHandlerInstance := supplierHandler.NewSupplierHandler(createSupplierUC, getSupplierUC, listSuppliersUC, updateSupplierUC, deleteSupplierUC)
	capitalInjectionHandlerInstance := capitalInjectionHandler.NewCapitalInjectionHandler(createInjectionUC, getInjectionUC, listInjectionsUC, getTotalCapitalUC)

	// Crear instancia de Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Configurar rutas
	routes.SetupRoutes(e, routes.Handlers{
		Auth:             authHandlerInstance,
		User:             userHandlerInstance,
		Product:          productHandlerInstance,
		Category:         categoryHandlerInstance,
		Size:             sizeHandlerInstance,
		PaymentMethod:    paymentMethodHandlerInstance,
		Customer:         customerHandlerInstance,
		Supplier:         supplierHandlerInstance,
		CapitalInjection: capitalInjectionHandlerInstance,
	}, validateTokenUC)

	// Iniciar servidor
	addr := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	log.Printf("Starting server on %s", addr)

	// Graceful shutdown
	go func() {
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Cerrar conexión a la base de datos
	if err := database.Close(db); err != nil {
		log.Fatal("Failed to close database connection:", err)
	}

	log.Println("Server stopped gracefully")
}
