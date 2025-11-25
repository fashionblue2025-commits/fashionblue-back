package product

import (
	"bytes"
	"context"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type UploadProductPhotoUseCase struct {
	productPhotoRepo ports.ProductPhotoRepository
	fileStorage      ports.FileStorage
}

func NewUploadProductPhotoUseCase(
	productPhotoRepo ports.ProductPhotoRepository,
	fileStorage ports.FileStorage,
) *UploadProductPhotoUseCase {
	return &UploadProductPhotoUseCase{
		productPhotoRepo: productPhotoRepo,
		fileStorage:      fileStorage,
	}
}

func (uc *UploadProductPhotoUseCase) Execute(
	ctx context.Context,
	productID uint,
	fileName string,
	fileData []byte,
	contentType string,
	description string,
	isPrimary bool,
) (*entities.ProductPhoto, error) {
	// Subir archivo
	reader := bytes.NewReader(fileData)
	photoURL, err := uc.fileStorage.UploadFile(ctx, reader, fileName, contentType)
	if err != nil {
		return nil, err
	}

	// Crear registro de foto
	photo := &entities.ProductPhoto{
		ProductID:   productID,
		PhotoURL:    photoURL,
		Description: description,
		IsPrimary:   isPrimary,
		UploadedAt:  time.Now(),
	}

	if err := uc.productPhotoRepo.Create(ctx, photo); err != nil {
		// Si falla la creación, intentar eliminar el archivo subido
		_ = uc.fileStorage.DeleteFile(ctx, photoURL)
		return nil, err
	}

	// Si es primary, actualizar las demás fotos
	if isPrimary {
		_ = uc.productPhotoRepo.SetAsPrimary(ctx, photo.ID, productID)
	}

	return photo, nil
}
