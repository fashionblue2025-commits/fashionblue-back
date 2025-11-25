package ports

import (
	"context"
	"io"
)

// FileStorage define la interfaz para almacenamiento de archivos
type FileStorage interface {
	// UploadFile sube un archivo y retorna la URL pública
	UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)

	// DeleteFile elimina un archivo por su URL
	DeleteFile(ctx context.Context, fileURL string) error

	// GetFileURL genera una URL firmada temporal (opcional)
	GetFileURL(ctx context.Context, filename string) (string, error)
}
