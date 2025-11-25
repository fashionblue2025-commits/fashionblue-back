package product

import (
	"bytes"
	"context"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type UploadMultiplePhotosUseCase struct {
	productPhotoRepo ports.ProductPhotoRepository
	fileStorage      ports.FileStorage
}

func NewUploadMultiplePhotosUseCase(
	productPhotoRepo ports.ProductPhotoRepository,
	fileStorage ports.FileStorage,
) *UploadMultiplePhotosUseCase {
	return &UploadMultiplePhotosUseCase{
		productPhotoRepo: productPhotoRepo,
		fileStorage:      fileStorage,
	}
}

// PhotoUpload representa los datos de una foto a subir
type PhotoUpload struct {
	FileName    string
	FileData    []byte
	ContentType string
	Description string
	IsPrimary   bool
}

func (uc *UploadMultiplePhotosUseCase) Execute(
	ctx context.Context,
	productID uint,
	photos []PhotoUpload,
) ([]entities.ProductPhoto, error) {
	uploadedPhotos := make([]entities.ProductPhoto, 0, len(photos))

	for i, photoData := range photos {
		// Subir archivo
		reader := bytes.NewReader(photoData.FileData)
		photoURL, err := uc.fileStorage.UploadFile(ctx, reader, photoData.FileName, photoData.ContentType)
		if err != nil {
			// Si falla, intentar limpiar las fotos ya subidas
			uc.cleanupUploadedPhotos(ctx, uploadedPhotos)
			return nil, err
		}

		// Crear registro de foto
		photo := entities.ProductPhoto{
			ProductID:    productID,
			PhotoURL:     photoURL,
			Description:  photoData.Description,
			IsPrimary:    photoData.IsPrimary,
			DisplayOrder: i + 1, // Orden basado en la posición en el array
			UploadedAt:   time.Now(),
		}

		if err := uc.productPhotoRepo.Create(ctx, &photo); err != nil {
			// Si falla la creación, eliminar el archivo subido y limpiar
			_ = uc.fileStorage.DeleteFile(ctx, photoURL)
			uc.cleanupUploadedPhotos(ctx, uploadedPhotos)
			return nil, err
		}

		uploadedPhotos = append(uploadedPhotos, photo)

		// Si es primary, actualizar las demás fotos
		if photoData.IsPrimary {
			_ = uc.productPhotoRepo.SetAsPrimary(ctx, photo.ID, productID)
		}
	}

	return uploadedPhotos, nil
}

// cleanupUploadedPhotos limpia las fotos subidas en caso de error
func (uc *UploadMultiplePhotosUseCase) cleanupUploadedPhotos(ctx context.Context, photos []entities.ProductPhoto) {
	for _, photo := range photos {
		_ = uc.fileStorage.DeleteFile(ctx, photo.PhotoURL)
		_ = uc.productPhotoRepo.Delete(ctx, photo.ID)
	}
}
