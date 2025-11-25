package order

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetOrderPhotosUseCase struct {
	orderPhotoRepo ports.OrderPhotoRepository
}

func NewGetOrderPhotosUseCase(orderPhotoRepo ports.OrderPhotoRepository) *GetOrderPhotosUseCase {
	return &GetOrderPhotosUseCase{
		orderPhotoRepo: orderPhotoRepo,
	}
}

func (uc *GetOrderPhotosUseCase) Execute(ctx context.Context, orderID uint) ([]entities.OrderPhoto, error) {
	return uc.orderPhotoRepo.GetByOrderID(ctx, orderID)
}
