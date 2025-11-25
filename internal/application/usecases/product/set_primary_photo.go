package product

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type SetPrimaryPhotoUseCase struct {
	productPhotoRepo ports.ProductPhotoRepository
}

func NewSetPrimaryPhotoUseCase(productPhotoRepo ports.ProductPhotoRepository) *SetPrimaryPhotoUseCase {
	return &SetPrimaryPhotoUseCase{
		productPhotoRepo: productPhotoRepo,
	}
}

func (uc *SetPrimaryPhotoUseCase) Execute(ctx context.Context, photoID uint, productID uint) error {
	return uc.productPhotoRepo.SetAsPrimary(ctx, photoID, productID)
}
