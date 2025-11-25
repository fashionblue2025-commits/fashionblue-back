package product

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetProductPhotosUseCase struct {
	productPhotoRepo ports.ProductPhotoRepository
}

func NewGetProductPhotosUseCase(productPhotoRepo ports.ProductPhotoRepository) *GetProductPhotosUseCase {
	return &GetProductPhotosUseCase{
		productPhotoRepo: productPhotoRepo,
	}
}

func (uc *GetProductPhotosUseCase) Execute(ctx context.Context, productID uint) ([]entities.ProductPhoto, error) {
	return uc.productPhotoRepo.GetByProductID(ctx, productID)
}
