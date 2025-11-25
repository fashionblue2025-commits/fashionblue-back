package database

import (
	"fmt"
	"log"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresDB crea una nueva conexión a PostgreSQL
func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	var logLevel logger.LogLevel
	if cfg.IsDevelopment() {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configurar pool de conexiones
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully")

	return db, nil
}

// AutoMigrate ejecuta las migraciones automáticas de GORM
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.UserModel{},
		&models.CategoryModel{},
		&models.ProductModel{},
		&models.ProductPhotoModel{},  // Tabla de fotos de productos
		&models.SizeModel{},          // Tabla de tallas
		&models.PaymentMethodModel{}, // Tabla de métodos de pago
		&models.CustomerModel{},
		&models.CustomerTransactionModel{},
		&models.SupplierModel{},   // Tabla de proveedores
		&models.OrderModel{},      // Tabla de órdenes
		&models.OrderItemModel{},  // Tabla de items de órdenes
		&models.OrderPhotoModel{}, // Tabla de fotos de órdenes
	)
}

// Close cierra la conexión a la base de datos
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
