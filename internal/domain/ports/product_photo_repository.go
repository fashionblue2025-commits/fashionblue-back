package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// ProductPhotoRepository define las operaciones para gestionar fotos de productos
type ProductPhotoRepository interface {
	Create(ctx context.Context, photo *entities.ProductPhoto) error
	GetByID(ctx context.Context, id uint) (*entities.ProductPhoto, error)
	GetByProductID(ctx context.Context, productID uint) ([]entities.ProductPhoto, error)
	Update(ctx context.Context, photo *entities.ProductPhoto) error
	Delete(ctx context.Context, id uint) error
	SetAsPrimary(ctx context.Context, photoID uint, productID uint) error // Establece una foto como principal
}
