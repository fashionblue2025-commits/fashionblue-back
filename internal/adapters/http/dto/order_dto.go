package dto

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// OrderDTO representa una orden en la API
type OrderDTO struct {
	ID                    uint            `json:"id"`
	OrderNumber           string          `json:"orderNumber"`
	CustomerID            *uint           `json:"customerId,omitempty"` // ID del cliente interno (opcional)
	CustomerName          string          `json:"customerName"`
	SellerID              uint            `json:"sellerId"`
	SellerName            string          `json:"sellerName,omitempty"`
	Type                  string          `json:"type"`
	Status                string          `json:"status"`
	TotalAmount           float64         `json:"totalAmount"`
	Discount              float64         `json:"discount"`
	Notes                 string          `json:"notes,omitempty"`
	OrderDate             time.Time       `json:"orderDate"`
	EstimatedDeliveryDate *time.Time      `json:"estimatedDeliveryDate,omitempty"`
	ActualDeliveryDate    *time.Time      `json:"actualDeliveryDate,omitempty"`
	Items                 []OrderItemDTO  `json:"items,omitempty"`
	Photos                []OrderPhotoDTO `json:"photos,omitempty"`
	CreatedAt             time.Time       `json:"createdAt"`
	UpdatedAt             time.Time       `json:"updatedAt"`
}

// OrderItemDTO representa un item de orden en la API
type OrderItemDTO struct {
	ID               uint        `json:"id"`
	OrderID          uint        `json:"orderId"`
	ProductID        uint        `json:"productId"`
	ProductName      string      `json:"productName"`
	CategoryID       uint        `json:"categoryId"`
	Product          *ProductDTO `json:"product,omitempty"`
	Color            string      `json:"color,omitempty"`
	SizeID           *uint       `json:"sizeId,omitempty"`
	SizeName         string      `json:"sizeName,omitempty"`
	Size             *SizeDTO    `json:"size,omitempty"`
	Quantity         int         `json:"quantity"`
	UnitPrice        float64     `json:"unitPrice"`
	Subtotal         float64     `json:"subtotal"`
	ReservedQuantity int         `json:"reservedQuantity"`
}

// OrderPhotoDTO representa una foto de orden en la API
type OrderPhotoDTO struct {
	ID          uint      `json:"id"`
	OrderID     uint      `json:"orderId"`
	PhotoURL    string    `json:"photoUrl"`
	Description string    `json:"description,omitempty"`
	UploadedAt  time.Time `json:"uploadedAt"`
}

// ToOrderDTO convierte una entidad Order a DTO
func ToOrderDTO(order *entities.Order) *OrderDTO {
	dto := &OrderDTO{
		ID:                    order.ID,
		OrderNumber:           order.OrderNumber,
		CustomerID:            order.CustomerID,
		CustomerName:          order.CustomerName,
		SellerID:              order.SellerID,
		Type:                  string(order.Type),
		Status:                string(order.Status),
		TotalAmount:           order.TotalAmount,
		Discount:              order.Discount,
		Notes:                 order.Notes,
		OrderDate:             order.OrderDate,
		EstimatedDeliveryDate: order.EstimatedDeliveryDate,
		ActualDeliveryDate:    order.ActualDeliveryDate,
		CreatedAt:             order.CreatedAt,
		UpdatedAt:             order.UpdatedAt,
	}

	// Agregar nombre del vendedor si existe
	if order.Seller != nil {
		dto.SellerName = order.Seller.FirstName + " " + order.Seller.LastName
	}

	// Convertir items
	if len(order.Items) > 0 {
		dto.Items = make([]OrderItemDTO, len(order.Items))
		for i, item := range order.Items {
			dto.Items[i] = *ToOrderItemDTO(&item)
		}
	}

	// Convertir fotos
	if len(order.Photos) > 0 {
		dto.Photos = make([]OrderPhotoDTO, len(order.Photos))
		for i, photo := range order.Photos {
			dto.Photos[i] = *ToOrderPhotoDTO(&photo)
		}
	}

	return dto
}

// ToOrderDTOList convierte un slice de órdenes a DTOs
func ToOrderDTOList(orders []entities.Order) []*OrderDTO {
	dtos := make([]*OrderDTO, len(orders))
	for i, order := range orders {
		dtos[i] = ToOrderDTO(&order)
	}
	return dtos
}

// ToOrderItemDTO convierte una entidad OrderItem a DTO
func ToOrderItemDTO(item *entities.OrderItem) *OrderItemDTO {
	dto := &OrderItemDTO{
		ID:               item.ID,
		OrderID:          item.OrderID,
		ProductID:        item.ProductVariantID, // Usar ProductVariantID
		ProductName:      item.ProductName,
		CategoryID:       item.CategoryID,
		Color:            item.Color,
		SizeID:           item.SizeID,
		Quantity:         item.Quantity,
		UnitPrice:        item.UnitPrice,
		Subtotal:         item.Subtotal,
		ReservedQuantity: item.ReservedQuantity,
	}

	// Agregar variante completa si existe
	if item.ProductVariant != nil && item.ProductVariant.Product != nil {
		dto.Product = ToProductDTO(item.ProductVariant.Product)
	}

	// Agregar talla completa si existe
	if item.Size != nil {
		dto.SizeName = item.Size.Value
		dto.Size = ToSizeDTO(item.Size)
	}

	return dto
}

// ToOrderPhotoDTO convierte una entidad OrderPhoto a DTO
func ToOrderPhotoDTO(photo *entities.OrderPhoto) *OrderPhotoDTO {
	return &OrderPhotoDTO{
		ID:          photo.ID,
		OrderID:     photo.OrderID,
		PhotoURL:    photo.PhotoURL,
		Description: photo.Description,
		UploadedAt:  photo.UploadedAt,
	}
}
