package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type LocalFileStorage struct {
	uploadDir string
	baseURL   string
}

func NewLocalFileStorage(uploadDir, baseURL string) ports.FileStorage {
	// Crear directorio si no existe
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create upload directory: %v", err))
	}

	return &LocalFileStorage{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

func (s *LocalFileStorage) UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	// Generar nombre único
	ext := filepath.Ext(filename)
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	uniqueFilename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), hex.EncodeToString(randomBytes), ext)

	// Crear subdirectorio por fecha
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(s.uploadDir, dateDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Ruta completa del archivo
	filePath := filepath.Join(fullDir, uniqueFilename)

	// Crear archivo
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copiar contenido
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Generar URL pública
	publicURL := fmt.Sprintf("%s/uploads/%s/%s", s.baseURL, dateDir, uniqueFilename)

	return publicURL, nil
}

func (s *LocalFileStorage) DeleteFile(ctx context.Context, fileURL string) error {
	// Extraer ruta del archivo desde la URL
	// Ejemplo: http://localhost:8080/uploads/2024/01/15/file.jpg -> uploads/2024/01/15/file.jpg
	// Implementación simplificada
	return nil
}

func (s *LocalFileStorage) GetFileURL(ctx context.Context, filename string) (string, error) {
	return fmt.Sprintf("%s/uploads/%s", s.baseURL, filename), nil
}
