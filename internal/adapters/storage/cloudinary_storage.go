package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryStorage struct {
	cld    *cloudinary.Cloudinary
	folder string
}

// NewCloudinaryStorage crea una nueva instancia de CloudinaryStorage
// cloudName: Tu cloud name de Cloudinary
// apiKey: Tu API key
// apiSecret: Tu API secret
// folder: Carpeta donde se guardarán las imágenes (ej: "fashion-blue/orders")
func NewCloudinaryStorage(cloudName, apiKey, apiSecret, folder string) (ports.FileStorage, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	return &CloudinaryStorage{
		cld:    cld,
		folder: folder,
	}, nil
}

func (s *CloudinaryStorage) UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	// Generar public ID único
	ext := filepath.Ext(filename)
	nameWithoutExt := filename[:len(filename)-len(ext)]
	publicID := fmt.Sprintf("%s_%d", nameWithoutExt, time.Now().Unix())

	// Configurar opciones de subida
	uniqueFilename := true
	overwrite := false
	uploadParams := uploader.UploadParams{
		PublicID:       publicID,
		Folder:         s.folder,
		ResourceType:   "image",
		UniqueFilename: &uniqueFilename,
		Overwrite:      &overwrite,
		Transformation: "q_auto,f_auto",
	}

	// Subir archivo
	result, err := s.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload to cloudinary: %w", err)
	}

	// Retornar URL segura (HTTPS)
	return result.SecureURL, nil
}

func (s *CloudinaryStorage) DeleteFile(ctx context.Context, fileURL string) error {
	// Extraer public ID de la URL
	// Cloudinary URL format: https://res.cloudinary.com/{cloud_name}/image/upload/{transformations}/{public_id}.{format}
	// Por simplicidad, se puede implementar después si es necesario
	// Por ahora retornamos nil (las imágenes se pueden gestionar desde el dashboard de Cloudinary)
	return nil
}

func (s *CloudinaryStorage) GetFileURL(ctx context.Context, filename string) (string, error) {
	// Cloudinary maneja las URLs automáticamente al subir
	// Este método puede retornar una URL con transformaciones específicas si se necesita
	publicID := fmt.Sprintf("%s/%s", s.folder, filename)
	url, _ := s.cld.Image(publicID)
	urlString, err := url.String()
	if err != nil {
		return "", fmt.Errorf("failed to generate cloudinary url: %w", err)
	}
	return urlString, nil
}
