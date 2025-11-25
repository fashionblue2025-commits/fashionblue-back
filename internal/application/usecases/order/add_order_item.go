package order

import (
	"context"
	"errors"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type AddOrderItemUseCase struct {
	orderRepo          ports.OrderRepository
	orderItemRepo      ports.OrderItemRepository
	productRepo        ports.ProductRepository
	productVariantRepo ports.ProductVariantRepository
}

func NewAddOrderItemUseCase(
	orderRepo ports.OrderRepository,
	orderItemRepo ports.OrderItemRepository,
	productRepo ports.ProductRepository,
	productVariantRepo ports.ProductVariantRepository,
) *AddOrderItemUseCase {
	return &AddOrderItemUseCase{
		orderRepo:          orderRepo,
		orderItemRepo:      orderItemRepo,
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
	}
}

func (uc *AddOrderItemUseCase) Execute(ctx context.Context, item *entities.OrderItem) error {
	// Validar item
	if err := item.Validate(); err != nil {
		return err
	}

	// Obtener orden
	order, err := uc.orderRepo.GetByID(ctx, item.OrderID)
	if err != nil {
		return err
	}

	// Verificar que se pueden editar items
	if !order.CanEditItems() {
		return errors.New("cannot edit items in current order status")
	}

	// Obtener variante para snapshot del nombre
	if item.ProductVariantID != 0 {
		variant, err := uc.productVariantRepo.GetByID(ctx, item.ProductVariantID)
		if err != nil {
			return err
		}
		if variant.Product != nil {
			item.ProductName = variant.Product.Name
		}
	}

	// Calcular subtotal
	item.CalculateSubtotal()

	// Crear item
	if err := uc.orderItemRepo.Create(ctx, item); err != nil {
		return err
	}

	// Recalcular total de la orden
	order.Items = append(order.Items, *item)
	order.TotalAmount = order.CalculateTotal()

	return uc.orderRepo.Update(ctx, order)
}
