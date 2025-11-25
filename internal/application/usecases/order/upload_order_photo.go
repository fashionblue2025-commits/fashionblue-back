package order

import (
	"context"
	"io"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type UploadOrderPhotoUseCase struct {
	orderPhotoRepo ports.OrderPhotoRepository
	fileStorage    ports.FileStorage
}

func NewUploadOrderPhotoUseCase(orderPhotoRepo ports.OrderPhotoRepository, fileStorage ports.FileStorage) *UploadOrderPhotoUseCase {
	return &UploadOrderPhotoUseCase{
		orderPhotoRepo: orderPhotoRepo,
		fileStorage:    fileStorage,
	}
}

// Execute sube el archivo a la nube y guarda la referencia en la base de datos
func (uc *UploadOrderPhotoUseCase) Execute(ctx context.Context, orderID uint, file io.Reader, filename string, contentType string, description string) (*entities.OrderPhoto, error) {
	// Subir archivo al almacenamiento
	photoURL, err := uc.fileStorage.UploadFile(ctx, file, filename, contentType)
	if err != nil {
		return nil, err
	}

	// Crear entidad de foto
	photo := &entities.OrderPhoto{
		OrderID:     orderID,
		PhotoURL:    photoURL,
		Description: description,
		UploadedAt:  time.Now(),
	}

	// Validar foto
	if err := photo.Validate(); err != nil {
		return nil, err
	}

	// Guardar en base de datos
	if err := uc.orderPhotoRepo.Create(ctx, photo); err != nil {
		// TODO: Eliminar archivo del almacenamiento si falla el guardado en BD
		return nil, err
	}

	return photo, nil
}
