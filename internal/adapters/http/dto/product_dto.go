package dto

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// ProductDTO representa un producto base en la API
type ProductDTO struct {
	ID              uint                `json:"id"`
	Name            string              `json:"name"`
	Description     string              `json:"description,omitempty"`
	CategoryID      uint                `json:"categoryId"`
	Category        *CategoryDTO        `json:"category,omitempty"`
	MaterialCost    float64             `json:"materialCost"`
	LaborCost       float64             `json:"laborCost"`
	ProductionCost  float64             `json:"productionCost"`
	UnitPrice       float64             `json:"unitPrice"`
	WholesalePrice  float64             `json:"wholesalePrice"`
	MinWholesaleQty int                 `json:"minWholesaleQty"`
	MinStock        int                 `json:"minStock"`
	IsActive        bool                `json:"isActive"`
	Variants        []ProductVariantDTO `json:"variants,omitempty"` // Variantes del producto
	TotalStock      int                 `json:"totalStock"`         // Stock total de todas las variantes
	CreatedAt       time.Time           `json:"createdAt"`
	UpdatedAt       time.Time           `json:"updatedAt"`
}

// ProductVariantDTO representa una variante de producto (color + talla)
type ProductVariantDTO struct {
	ID             uint      `json:"id"`
	ProductID      uint      `json:"productId"`
	Color          string    `json:"color"`
	SizeID         *uint     `json:"sizeId,omitempty"`
	Size           *SizeDTO  `json:"size,omitempty"`
	Stock          int       `json:"stock"`
	ReservedStock  int       `json:"reservedStock"`
	AvailableStock int       `json:"availableStock"`
	UnitPrice      float64   `json:"unitPrice"`
	IsActive       bool      `json:"isActive"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// CategoryDTO representa una categoría en la API
type CategoryDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"` // Identificador para tipo de tallas
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// SizeDTO representa una talla en la API
type SizeDTO struct {
	ID        uint      `json:"id"`
	Value     string    `json:"value"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToProductDTO convierte una entidad Product a DTO
func ToProductDTO(product *entities.Product) *ProductDTO {
	dto := &ProductDTO{
		ID:              product.ID,
		Name:            product.Name,
		Description:     product.Description,
		CategoryID:      product.CategoryID,
		MaterialCost:    product.MaterialCost,
		LaborCost:       product.LaborCost,
		ProductionCost:  product.ProductionCost,
		UnitPrice:       product.UnitPrice,
		WholesalePrice:  product.WholesalePrice,
		MinWholesaleQty: product.MinWholesaleQty,
		MinStock:        product.MinStock,
		IsActive:        product.IsActive,
		TotalStock:      product.GetTotalStock(),
		CreatedAt:       product.CreatedAt,
		UpdatedAt:       product.UpdatedAt,
	}

	// Convertir Category si existe
	if product.Category != nil {
		dto.Category = ToCategoryDTO(product.Category)
	}

	// Convertir Variants si existen
	if len(product.Variants) > 0 {
		dto.Variants = make([]ProductVariantDTO, len(product.Variants))
		for i, variant := range product.Variants {
			dto.Variants[i] = *ToProductVariantDTO(&variant)
		}
	}

	return dto
}

// ToProductVariantDTO convierte una entidad ProductVariant a DTO
func ToProductVariantDTO(variant *entities.ProductVariant) *ProductVariantDTO {
	dto := &ProductVariantDTO{
		ID:             variant.ID,
		ProductID:      variant.ProductID,
		Color:          variant.Color,
		SizeID:         variant.SizeID,
		Stock:          variant.Stock,
		ReservedStock:  variant.ReservedStock,
		AvailableStock: variant.GetAvailableStock(),
		UnitPrice:      variant.UnitPrice,
		IsActive:       variant.IsActive,
		CreatedAt:      variant.CreatedAt,
		UpdatedAt:      variant.UpdatedAt,
	}

	// Convertir Size si existe
	if variant.Size != nil {
		dto.Size = ToSizeDTO(variant.Size)
	}

	return dto
}

// ToProductDTOList convierte un slice de productos a DTOs
func ToProductDTOList(products []entities.Product) []*ProductDTO {
	dtos := make([]*ProductDTO, len(products))
	for i, product := range products {
		dtos[i] = ToProductDTO(&product)
	}
	return dtos
}

// ToCategoryDTO convierte una entidad Category a DTO
func ToCategoryDTO(category *entities.Category) *CategoryDTO {
	return &CategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

// ToCategoryDTOList convierte un slice de categorías a DTOs
func ToCategoryDTOList(categories []entities.Category) []*CategoryDTO {
	dtos := make([]*CategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = ToCategoryDTO(&category)
	}
	return dtos
}

// ToSizeDTO convierte una entidad Size a DTO
func ToSizeDTO(size *entities.Size) *SizeDTO {
	return &SizeDTO{
		ID:        size.ID,
		Value:     size.Value,
		Type:      string(size.Type),
		CreatedAt: size.CreatedAt,
		UpdatedAt: size.UpdatedAt,
	}
}

// ToSizeDTOList convierte un slice de tallas a DTOs
func ToSizeDTOList(sizes []entities.Size) []*SizeDTO {
	dtos := make([]*SizeDTO, len(sizes))
	for i, size := range sizes {
		dtos[i] = ToSizeDTO(&size)
	}
	return dtos
}

// ProductPhotoDTO representa una foto de producto en la API
type ProductPhotoDTO struct {
	ID           uint      `json:"id"`
	ProductID    uint      `json:"productId"`
	PhotoURL     string    `json:"photoUrl"`
	Description  string    `json:"description,omitempty"`
	IsPrimary    bool      `json:"isPrimary"`
	DisplayOrder int       `json:"displayOrder"`
	UploadedAt   time.Time `json:"uploadedAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// UploadPhotoRequestDTO representa la solicitud para subir una foto
type UploadPhotoRequestDTO struct {
	Description  string `json:"description,omitempty" form:"description"`
	IsPrimary    bool   `json:"isPrimary" form:"isPrimary"`
	DisplayOrder int    `json:"displayOrder" form:"displayOrder"`
}

// SetPrimaryPhotoRequestDTO representa la solicitud para establecer foto principal
type SetPrimaryPhotoRequestDTO struct {
	PhotoID uint `json:"photoId" validate:"required"`
}

// ToProductPhotoDTO convierte una entidad ProductPhoto a DTO
func ToProductPhotoDTO(photo *entities.ProductPhoto) *ProductPhotoDTO {
	return &ProductPhotoDTO{
		ID:           photo.ID,
		ProductID:    photo.ProductID,
		PhotoURL:     photo.PhotoURL,
		Description:  photo.Description,
		IsPrimary:    photo.IsPrimary,
		DisplayOrder: photo.DisplayOrder,
		UploadedAt:   photo.UploadedAt,
		CreatedAt:    photo.CreatedAt,
		UpdatedAt:    photo.UpdatedAt,
	}
}

// ToProductPhotoDTOList convierte un slice de fotos a DTOs
func ToProductPhotoDTOList(photos []entities.ProductPhoto) []*ProductPhotoDTO {
	dtos := make([]*ProductPhotoDTO, len(photos))
	for i, photo := range photos {
		dtos[i] = ToProductPhotoDTO(&photo)
	}
	return dtos
}
