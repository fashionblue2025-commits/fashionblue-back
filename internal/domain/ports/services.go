package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// AuthService define las operaciones de autenticación
type AuthService interface {
	Login(ctx context.Context, email, password string) (string, *entities.User, error)
	ValidateToken(ctx context.Context, token string) (*entities.User, error)
	GenerateToken(user *entities.User) (string, error)
}

// UserService define las operaciones de negocio para usuarios
type UserService interface {
	CreateUser(ctx context.Context, user *entities.User, password string) error
	GetUser(ctx context.Context, id uint) (*entities.User, error)
	ListUsers(ctx context.Context, filters map[string]interface{}) ([]entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id uint) error
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
}

// CategoryService define las operaciones de negocio para categorías
type CategoryService interface {
	CreateCategory(ctx context.Context, category *entities.Category) error
	GetCategory(ctx context.Context, id uint) (*entities.Category, error)
	ListCategories(ctx context.Context, filters map[string]interface{}) ([]entities.Category, error)
	UpdateCategory(ctx context.Context, category *entities.Category) error
	DeleteCategory(ctx context.Context, id uint) error
}

// ProductService define las operaciones de negocio para productos
type ProductService interface {
	CreateProduct(ctx context.Context, product *entities.Product) error
	GetProduct(ctx context.Context, id uint) (*entities.Product, error)
	ListProducts(ctx context.Context, filters map[string]interface{}) ([]entities.Product, error)
	ListProductsByCategory(ctx context.Context, categoryID uint) ([]entities.Product, error)
	UpdateProduct(ctx context.Context, product *entities.Product) error
	DeleteProduct(ctx context.Context, id uint) error
	GetLowStockProducts(ctx context.Context) ([]entities.Product, error)
}

// CustomerService define las operaciones de negocio para clientes
type CustomerService interface {
	CreateCustomer(ctx context.Context, customer *entities.Customer) error
	GetCustomer(ctx context.Context, id uint) (*entities.Customer, error)
	ListCustomers(ctx context.Context, filters map[string]interface{}) ([]entities.Customer, error)
	UpdateCustomer(ctx context.Context, customer *entities.Customer) error
	DeleteCustomer(ctx context.Context, id uint) error
	GetCustomerHistory(ctx context.Context, customerID uint) ([]entities.CustomerTransaction, error)
}

// SupplierService define las operaciones de negocio para proveedores
type SupplierService interface {
	CreateSupplier(ctx context.Context, supplier *entities.Supplier) error
	GetSupplier(ctx context.Context, id uint) (*entities.Supplier, error)
	ListSuppliers(ctx context.Context, filters map[string]interface{}) ([]entities.Supplier, error)
	UpdateSupplier(ctx context.Context, supplier *entities.Supplier) error
	DeleteSupplier(ctx context.Context, id uint) error
}

// FileService define las operaciones para manejo de archivos
type FileService interface {
	SaveFile(ctx context.Context, file []byte, filename string) (string, error)
	DeleteFile(ctx context.Context, filePath string) error
	GetFile(ctx context.Context, filePath string) ([]byte, error)
}
