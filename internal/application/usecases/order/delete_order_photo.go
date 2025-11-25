package order

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type DeleteOrderPhotoUseCase struct {
	orderPhotoRepo ports.OrderPhotoRepository
}

func NewDeleteOrderPhotoUseCase(orderPhotoRepo ports.OrderPhotoRepository) *DeleteOrderPhotoUseCase {
	return &DeleteOrderPhotoUseCase{
		orderPhotoRepo: orderPhotoRepo,
	}
}

func (uc *DeleteOrderPhotoUseCase) Execute(ctx context.Context, photoID uint) error {
	return uc.orderPhotoRepo.Delete(ctx, photoID)
}
