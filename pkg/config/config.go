package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	App        AppConfig
	Database   DatabaseConfig
	JWT        JWTConfig
	Upload     UploadConfig
	Cloudinary CloudinaryConfig
	CORS       CORSConfig
	Log        LogConfig
}

// AppConfig configuración de la aplicación
type AppConfig struct {
	Name string
	Env  string
	Port string
	Host string
}

// DatabaseConfig configuración de la base de datos
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig configuración de JWT
type JWTConfig struct {
	Secret     string
	Expiration string
}

// GetExpiration convierte la expiración de string a time.Duration
func (j *JWTConfig) GetExpiration() time.Duration {
	duration, err := time.ParseDuration(j.Expiration)
	if err != nil {
		return 24 * time.Hour // Default 24 horas
	}
	return duration
}

// UploadConfig configuración de uploads
type UploadConfig struct {
	MaxSize int64
	Path    string
}

// CloudinaryConfig configuración de Cloudinary
type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
	Folder    string
	Enabled   bool
}

// CORSConfig configuración de CORS
type CORSConfig struct {
	AllowedOrigins string
}

// LogConfig configuración de logs
type LogConfig struct {
	Level  string
	Format string
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	// Cargar archivo .env si existe
	_ = godotenv.Load()

	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)

	config := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "fashion-blue"),
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
			Host: getEnv("APP_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "fashionblue"),
			Password: getEnv("DB_PASSWORD", "fashionblue123"),
			Name:     getEnv("DB_NAME", "fashionblue_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
			Expiration: getEnv("JWT_EXPIRATION", "24h"),
		},
		Upload: UploadConfig{
			MaxSize: maxUploadSize,
			Path:    getEnv("UPLOAD_PATH", "./uploads"),
		},
		Cloudinary: CloudinaryConfig{
			CloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
			APIKey:    getEnv("CLOUDINARY_API_KEY", ""),
			APISecret: getEnv("CLOUDINARY_API_SECRET", ""),
			Folder:    getEnv("CLOUDINARY_FOLDER", "fashion-blue/orders"),
			Enabled:   getEnv("CLOUDINARY_ENABLED", "false") == "true",
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	return config, nil
}

// GetDSN retorna el DSN de conexión a la base de datos
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// IsDevelopment verifica si el entorno es desarrollo
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsProduction verifica si el entorno es producción
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
