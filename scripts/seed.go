package main

import (
	"fmt"
	"log"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/config"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

	// Ejecutar migraciones
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Crear datos iniciales
	if err := seedData(db); err != nil {
		log.Fatal("Failed to seed data:", err)
	}

	log.Println("Database seeded successfully!")
}

func seedData(db *gorm.DB) error {
	fmt.Println("\n🌱 Seeding database...")
	fmt.Println("=" + string(make([]byte, 50)))

	// 1. Crear usuario super admin
	fmt.Println("\n👤 Creating Super Admin...")
	var userCount int64
	db.Model(&models.UserModel{}).Where("role = ?", entities.RoleSuperAdmin).Count(&userCount)

	if userCount == 0 {
		// Hash de la contraseña
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		adminUser := &models.UserModel{
			Email:     "admin@fashionblue.com",
			Password:  string(hashedPassword),
			FirstName: "Super",
			LastName:  "Admin",
			Role:      string(entities.RoleSuperAdmin),
			IsActive:  true,
		}

		if err := db.Create(adminUser).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		fmt.Println("✅ Super Admin created:")
		fmt.Println("   Email: admin@fashionblue.com")
		fmt.Println("   Password: admin123")
	} else {
		fmt.Println("⏭️  Super Admin already exists")
	}

	// 2. Crear categorías iniciales
	fmt.Println("\n📁 Creating Categories...")
	categories := []models.CategoryModel{
		{Name: "Chaquetas", Description: "Chaquetas y abrigos de cuero", IsActive: true},
		{Name: "Pantalones", Description: "Pantalones y jeans de cuero", IsActive: true},
		{Name: "Camisas", Description: "Camisas y blusas", IsActive: true},
		{Name: "Accesorios", Description: "Cinturones, carteras y más", IsActive: true},
		{Name: "Calzado", Description: "Zapatos y botas de cuero", IsActive: true},
	}

	for _, category := range categories {
		var existingCategory models.CategoryModel
		result := db.Where("name = ?", category.Name).First(&existingCategory)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&category).Error; err != nil {
				log.Printf("❌ Error creating category '%s': %v", category.Name, err)
			} else {
				fmt.Printf("✅ Category created: %s\n", category.Name)
			}
		} else {
			fmt.Printf("⏭️  Category already exists: %s\n", category.Name)
		}
	}

	// 3. Crear tallas
	fmt.Println("\n📏 Creating Sizes...")
	shirtSizes := []models.SizeModel{
		{Type: entities.SizeTypeShirt, Value: "XS", Order: 1, IsActive: true},
		{Type: entities.SizeTypeShirt, Value: "S", Order: 2, IsActive: true},
		{Type: entities.SizeTypeShirt, Value: "M", Order: 3, IsActive: true},
		{Type: entities.SizeTypeShirt, Value: "L", Order: 4, IsActive: true},
		{Type: entities.SizeTypeShirt, Value: "XL", Order: 5, IsActive: true},
		{Type: entities.SizeTypeShirt, Value: "XXL", Order: 6, IsActive: true},
	}

	// Tallas de pantalones (cintura en pulgadas)
	pantsSizes := []models.SizeModel{
		{Type: entities.SizeTypePants, Value: "24", Order: 1, IsActive: true},
		{Type: entities.SizeTypePants, Value: "26", Order: 2, IsActive: true},
		{Type: entities.SizeTypePants, Value: "28", Order: 3, IsActive: true},
		{Type: entities.SizeTypePants, Value: "30", Order: 4, IsActive: true},
		{Type: entities.SizeTypePants, Value: "32", Order: 5, IsActive: true},
		{Type: entities.SizeTypePants, Value: "34", Order: 6, IsActive: true},
		{Type: entities.SizeTypePants, Value: "36", Order: 7, IsActive: true},
		{Type: entities.SizeTypePants, Value: "38", Order: 8, IsActive: true},
		{Type: entities.SizeTypePants, Value: "40", Order: 9, IsActive: true},
		{Type: entities.SizeTypePants, Value: "42", Order: 10, IsActive: true},
	}

	// Tallas de zapatos/tenis (sistema US)
	shoesSizes := []models.SizeModel{
		{Type: entities.SizeTypeShoes, Value: "5", Order: 1, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "5.5", Order: 2, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "6", Order: 3, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "6.5", Order: 4, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "7", Order: 5, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "7.5", Order: 6, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "8", Order: 7, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "8.5", Order: 8, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "9", Order: 9, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "9.5", Order: 10, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "10", Order: 11, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "10.5", Order: 12, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "11", Order: 13, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "11.5", Order: 14, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "12", Order: 15, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "13", Order: 16, IsActive: true},
		{Type: entities.SizeTypeShoes, Value: "14", Order: 17, IsActive: true},
	}

	// Insertar tallas de camisetas
	shirtCount := 0
	for _, size := range shirtSizes {
		result := db.FirstOrCreate(&size, models.SizeModel{Type: size.Type, Value: size.Value})
		if result.Error != nil {
			log.Printf("❌ Error creating shirt size %s: %v", size.Value, result.Error)
		} else if result.RowsAffected > 0 {
			fmt.Printf("✅ Shirt size created: %s\n", size.Value)
			shirtCount++
		}
	}
	if shirtCount == 0 {
		fmt.Println("⏭️  All shirt sizes already exist")
	}

	// Insertar tallas de pantalones
	pantsCount := 0
	for _, size := range pantsSizes {
		result := db.FirstOrCreate(&size, models.SizeModel{Type: size.Type, Value: size.Value})
		if result.Error != nil {
			log.Printf("❌ Error creating pants size %s: %v", size.Value, result.Error)
		} else if result.RowsAffected > 0 {
			fmt.Printf("✅ Pants size created: %s\n", size.Value)
			pantsCount++
		}
	}
	if pantsCount == 0 {
		fmt.Println("⏭️  All pants sizes already exist")
	}

	// Insertar tallas de zapatos
	shoesCount := 0
	for _, size := range shoesSizes {
		result := db.FirstOrCreate(&size, models.SizeModel{Type: size.Type, Value: size.Value})
		if result.Error != nil {
			log.Printf("❌ Error creating shoes size %s: %v", size.Value, result.Error)
		} else if result.RowsAffected > 0 {
			fmt.Printf("✅ Shoes size created: %s\n", size.Value)
			shoesCount++
		}
	}
	if shoesCount == 0 {
		fmt.Println("⏭️  All shoes sizes already exist")
	}

	fmt.Println("\n" + string(make([]byte, 50)))
	fmt.Println("✅ DATABASE SEEDED SUCCESSFULLY!")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println("\n📊 Summary:")
	fmt.Printf("   👤 Users: 1 (Super Admin)\n")
	fmt.Printf("   📁 Categories: %d\n", len(categories))
	fmt.Printf("   📏 Sizes: %d total\n", len(shirtSizes)+len(pantsSizes)+len(shoesSizes))
	fmt.Printf("      - Shirts: %d\n", len(shirtSizes))
	fmt.Printf("      - Pants: %d\n", len(pantsSizes))
	fmt.Printf("      - Shoes: %d\n", len(shoesSizes))
	fmt.Println("\n🔐 Login credentials:")
	fmt.Println("   Email: admin@fashionblue.com")
	fmt.Println("   Password: admin123")

	return nil
}
