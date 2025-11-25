package product

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type DeleteProductPhotoUseCase struct {
	productPhotoRepo ports.ProductPhotoRepository
	fileStorage      ports.FileStorage
}

func NewDeleteProductPhotoUseCase(
	productPhotoRepo ports.ProductPhotoRepository,
	fileStorage ports.FileStorage,
) *DeleteProductPhotoUseCase {
	return &DeleteProductPhotoUseCase{
		productPhotoRepo: productPhotoRepo,
		fileStorage:      fileStorage,
	}
}

func (uc *DeleteProductPhotoUseCase) Execute(ctx context.Context, photoID uint) error {
	// Obtener la foto para tener la URL
	photo, err := uc.productPhotoRepo.GetByID(ctx, photoID)
	if err != nil {
		return err
	}

	// Eliminar el archivo
	if err := uc.fileStorage.DeleteFile(ctx, photo.PhotoURL); err != nil {
		// Log error pero continuar con la eliminación del registro
		// En producción, podrías querer manejar esto de otra manera
	}

	// Eliminar el registro
	return uc.productPhotoRepo.Delete(ctx, photoID)
}
