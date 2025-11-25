package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/routes"
	auditRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/audit"
	categoryRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/category"
	customerRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/customer"
	financialTransactionRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/financial_transaction"
	orderRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/order"
	paymentMethodRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/payment_method"
	productRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/product"
	sizeRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/size"
	supplierRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/supplier"
	userRepo "github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/repositories/user"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/storage"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/event_handlers"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/auth"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/category"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/customer"
	financialTransactionUseCases "github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/financial_transaction"
	orderUseCases "github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/order"
	paymentMethodUseCases "github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/payment_method"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/product"
	sizeUseCases "github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/size"
	supplierUseCases "github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/supplier"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/user"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
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
	productVariantRepository := productRepo.NewProductVariantRepository(db)
	productPhotoRepository := productRepo.NewProductPhotoRepository(db)
	sizeRepository := sizeRepo.NewSizeRepository(db)
	paymentMethodRepository := paymentMethodRepo.NewPaymentMethodRepository(db)
	customerRepository := customerRepo.NewCustomerRepository(db)
	customerTransactionRepository := customerRepo.NewCustomerTransactionRepository(db)
	auditLogRepository := auditRepo.NewAuditLogRepository(db)
	supplierRepository := supplierRepo.NewSupplierRepository(db)
	financialTransactionRepository := financialTransactionRepo.NewFinancialTransactionRepository(db)
	orderRepository := orderRepo.NewOrderRepository(db)
	orderItemRepository := orderRepo.NewOrderItemRepository(db)

	// Inicializar almacenamiento de archivos
	var fileStorage ports.FileStorage
	if cfg.Cloudinary.Enabled {
		// Usar Cloudinary si está habilitado
		cloudStorage, err := storage.NewCloudinaryStorage(
			cfg.Cloudinary.CloudName,
			cfg.Cloudinary.APIKey,
			cfg.Cloudinary.APISecret,
			cfg.Cloudinary.Folder,
		)
		if err != nil {
			log.Printf("Warning: Failed to initialize Cloudinary, falling back to local storage: %v", err)
			fileStorage = storage.NewLocalFileStorage(cfg.Upload.Path, fmt.Sprintf("http://%s:%s", cfg.App.Host, cfg.App.Port))
		} else {
			fileStorage = cloudStorage
			log.Println("Using Cloudinary for file storage")
		}
	} else {
		// Usar almacenamiento local por defecto
		fileStorage = storage.NewLocalFileStorage(cfg.Upload.Path, fmt.Sprintf("http://%s:%s", cfg.App.Host, cfg.App.Port))
		log.Println("Using local file storage")
	}

	// Inicializar casos de uso - Auth
	loginUC := auth.NewLoginUseCase(userRepository, cfg.JWT.Secret, cfg.JWT.GetExpiration())
	registerUC := auth.NewRegisterUseCase(userRepository, cfg.JWT.Secret, cfg.JWT.GetExpiration())
	validateTokenUC := auth.NewValidateTokenUseCase(userRepository, cfg.JWT.Secret)

	// Inicializar Event Bus
	eventBus := events.NewEventBus()

	// Inicializar Event Handlers
	loggingHandler := event_handlers.NewLoggingHandler(eventBus)
	loggingHandler.Start()

	notificationHandler := event_handlers.NewNotificationHandler(eventBus)
	notificationHandler.Start()

	analyticsEventHandler := event_handlers.NewAnalyticsHandler(eventBus)
	analyticsEventHandler.Start()

	auditEventHandler := event_handlers.NewAuditHandler(eventBus, auditLogRepository)
	auditEventHandler.Start()

	// Product creation handler para órdenes INVENTORY
	productCreationHandler := event_handlers.NewProductCreationHandler(eventBus, productRepository, productVariantRepository, orderItemRepository)
	productCreationHandler.Start()

	// Internal customer transaction handler para registro contable
	internalCustomerTransactionHandler := event_handlers.NewInternalCustomerTransactionHandler(eventBus, customerTransactionRepository)
	internalCustomerTransactionHandler.Start()

	// Financial income handler para registrar ingresos automáticos por ventas
	financialIncomeHandler := event_handlers.NewFinancialIncomeHandler(eventBus, financialTransactionRepository)
	financialIncomeHandler.Start()

	// Webhook handler (opcional - configurar según necesidad)
	webhookConfig := event_handlers.WebhookConfig{
		URL:     "", // Configurar URL si se necesita
		Enabled: false,
		Secret:  "",
	}
	webhookHandler := event_handlers.NewWebhookHandler(eventBus, webhookConfig)
	webhookHandler.Start()

	log.Println("✅ Event handlers initialized and started")

	// Inicializar casos de uso - User
	createUserUC := user.NewCreateUserUseCase(userRepository)
	getUserUC := user.NewGetUserUseCase(userRepository)
	listUsersUC := user.NewListUsersUseCase(userRepository)
	updateUserUC := user.NewUpdateUserUseCase(userRepository)
	deleteUserUC := user.NewDeleteUserUseCase(userRepository)
	changePasswordUC := user.NewChangePasswordUseCase(userRepository)

	// Inicializar casos de uso - Product
	createProductUC := product.NewCreateProductUseCase(productRepository, productVariantRepository)
	getProductUC := product.NewGetProductUseCase(productRepository)
	listProductsUC := product.NewListProductsUseCase(productRepository)
	updateProductUC := product.NewUpdateProductUseCase(productRepository)
	deleteProductUC := product.NewDeleteProductUseCase(productRepository)
	getLowStockUC := product.NewGetLowStockProductsUseCase(productRepository)
	uploadProductPhotoUC := product.NewUploadProductPhotoUseCase(productPhotoRepository, fileStorage)
	uploadMultiplePhotosUC := product.NewUploadMultiplePhotosUseCase(productPhotoRepository, fileStorage)
	getProductPhotosUC := product.NewGetProductPhotosUseCase(productPhotoRepository)
	deleteProductPhotoUC := product.NewDeleteProductPhotoUseCase(productPhotoRepository, fileStorage)
	setPrimaryPhotoUC := product.NewSetPrimaryPhotoUseCase(productPhotoRepository)

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
	generateCustomerStatementUC := usecases.NewGenerateCustomerStatementUseCase(customerRepository, customerTransactionRepository)

	// Inicializar casos de uso - Supplier
	createSupplierUC := supplierUseCases.NewCreateSupplierUseCase(supplierRepository)
	getSupplierUC := supplierUseCases.NewGetSupplierUseCase(supplierRepository)
	listSuppliersUC := supplierUseCases.NewListSuppliersUseCase(supplierRepository)
	updateSupplierUC := supplierUseCases.NewUpdateSupplierUseCase(supplierRepository)
	deleteSupplierUC := supplierUseCases.NewDeleteSupplierUseCase(supplierRepository)

	// Inicializar casos de uso - FinancialTransaction
	createTransactionUC := financialTransactionUseCases.NewCreateTransactionUseCase(financialTransactionRepository)
	getTransactionUC := financialTransactionUseCases.NewGetTransactionUseCase(financialTransactionRepository)
	listTransactionsUC := financialTransactionUseCases.NewListTransactionsUseCase(financialTransactionRepository)
	getBalanceUC := financialTransactionUseCases.NewGetBalanceUseCase(financialTransactionRepository)

	// Inicializar casos de uso - Order
	createOrderUC := orderUseCases.NewCreateOrderUseCase(orderRepository, productRepository, productVariantRepository, eventBus)
	getOrderUC := orderUseCases.NewGetOrderUseCase(orderRepository)
	listOrdersUC := orderUseCases.NewListOrdersUseCase(orderRepository)
	updateOrderStatusUC := orderUseCases.NewUpdateOrderStatusUseCase(orderRepository)
	addOrderItemUC := orderUseCases.NewAddOrderItemUseCase(orderRepository, orderItemRepository, productRepository, productVariantRepository)
	updateOrderItemUC := orderUseCases.NewUpdateOrderItemUseCase(orderRepository, orderItemRepository)
	removeOrderItemUC := orderUseCases.NewRemoveOrderItemUseCase(orderRepository, orderItemRepository)
	changeOrderStatusUC := orderUseCases.NewChangeOrderStatusUseCase(orderRepository, orderItemRepository, productRepository, productVariantRepository, eventBus)

	// Inicializar handlers
	authHandlerInstance := authHandler.NewAuthHandler(loginUC, registerUC)
	userHandlerInstance := userHandler.NewUserHandler(createUserUC, getUserUC, listUsersUC, updateUserUC, deleteUserUC, changePasswordUC)
	productHandlerInstance := productHandler.NewProductHandler(createProductUC, getProductUC, listProductsUC, updateProductUC, deleteProductUC, getLowStockUC, uploadProductPhotoUC, uploadMultiplePhotosUC, getProductPhotosUC, deleteProductPhotoUC, setPrimaryPhotoUC)
	categoryHandlerInstance := categoryHandler.NewCategoryHandler(createCategoryUC, getCategoryUC, listCategoriesUC, updateCategoryUC, deleteCategoryUC)
	sizeHandlerInstance := sizeHandler.NewSizeHandler(listSizesUC, getSizeUC, getSizesByTypeUC)
	paymentMethodHandlerInstance := paymentMethodHandler.NewPaymentMethodHandler(listPaymentMethodsUC)
	customerHandlerInstance := customerHandler.NewCustomerHandler(createCustomerUC, getCustomerUC, listCustomersUC, updateCustomerUC, deleteCustomerUC, getCustomerHistoryUC, createPaymentUC, getUpcomingPaymentsUC, getCustomerBalanceUC, addTransactionUC)
	statementHandlerInstance := customerHandler.NewStatementHandler(generateCustomerStatementUC)
	orderHandlerInstance := orderHandler.NewOrderHandler(createOrderUC, getOrderUC, listOrdersUC, updateOrderStatusUC, addOrderItemUC, updateOrderItemUC, removeOrderItemUC, changeOrderStatusUC)
	supplierHandlerInstance := supplierHandler.NewSupplierHandler(createSupplierUC, getSupplierUC, listSuppliersUC, updateSupplierUC, deleteSupplierUC)
	financialTransactionHandlerInstance := financialTransactionHandler.NewFinancialTransactionHandler(createTransactionUC, getTransactionUC, listTransactionsUC, getBalanceUC)
	analyticsHTTPHandlerInstance := analyticsHandler.NewAnalyticsHTTPHandler(analyticsEventHandler)
	auditHTTPHandlerInstance := auditHandler.NewAuditHTTPHandler(auditLogRepository)
	swaggerHandlerInstance := swaggerHandler.NewSwaggerHandler("docs/swagger.json")

	// Crear instancia de Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	// Servir archivos estáticos (uploads)
	e.Static("/uploads", "./uploads")

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Configurar rutas
	routes.SetupRoutes(e, routes.Handlers{
		Analytics:            analyticsHTTPHandlerInstance,
		Audit:                auditHTTPHandlerInstance,
		Auth:                 authHandlerInstance,
		User:                 userHandlerInstance,
		Product:              productHandlerInstance,
		Category:             categoryHandlerInstance,
		Size:                 sizeHandlerInstance,
		PaymentMethod:        paymentMethodHandlerInstance,
		Customer:             customerHandlerInstance,
		CustomerStatement:    statementHandlerInstance,
		Order:                orderHandlerInstance,
		Supplier:             supplierHandlerInstance,
		FinancialTransaction: financialTransactionHandlerInstance,
		Swagger:              swaggerHandlerInstance,
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

	// Detener event handlers
	log.Println("Stopping event handlers...")
	loggingHandler.Stop()
	notificationHandler.Stop()
	analyticsEventHandler.Stop()
	auditEventHandler.Stop()
	productCreationHandler.Stop()
	webhookHandler.Stop()

	// Cerrar event bus
	eventBus.Close()
	log.Println("Event handlers stopped")

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
