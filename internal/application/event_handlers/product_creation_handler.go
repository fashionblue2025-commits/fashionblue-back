package event_handlers

import (
	"context"
	"log"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

// ProductCreationHandler maneja la creación automática de productos y variantes
type ProductCreationHandler struct {
	eventBus           *events.EventBus
	eventChan          chan events.OrderEvent
	stopChan           chan bool
	productRepo        ports.ProductRepository
	productVariantRepo ports.ProductVariantRepository
	orderItemRepo      ports.OrderItemRepository
}

// NewProductCreationHandler crea un nuevo handler de creación de productos
func NewProductCreationHandler(
	eventBus *events.EventBus,
	productRepo ports.ProductRepository,
	productVariantRepo ports.ProductVariantRepository,
	orderItemRepo ports.OrderItemRepository,
) *ProductCreationHandler {
	handler := &ProductCreationHandler{
		eventBus:           eventBus,
		eventChan:          make(chan events.OrderEvent, 100),
		stopChan:           make(chan bool),
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
		orderItemRepo:      orderItemRepo,
	}

	// Suscribirse al evento de creación de productos
	eventBus.Subscribe(events.EventProductCreationRequired, handler.eventChan)

	return handler
}

// Start inicia el procesamiento de eventos
func (h *ProductCreationHandler) Start() {
	log.Println("🏭 Product Creation Event Handler started")

	go func() {
		for {
			select {
			case event := <-h.eventChan:
				h.handleEvent(event)
			case <-h.stopChan:
				log.Println("🏭 Product Creation Event Handler stopped")
				return
			}
		}
	}()
}

// Stop detiene el procesamiento de eventos
func (h *ProductCreationHandler) Stop() {
	h.stopChan <- true
}

// handleEvent procesa el evento de creación de productos
func (h *ProductCreationHandler) handleEvent(event events.OrderEvent) {
	if event.Type != events.EventProductCreationRequired {
		return
	}

	ctx := context.Background()

	// Obtener tipo de orden
	orderType, ok := event.Data["orderType"].(entities.OrderType)
	if !ok {
		log.Printf("🏭 [ERROR] Invalid orderType in event data")
		return
	}

	// Solo procesar órdenes INVENTORY y CUSTOM
	if orderType != entities.OrderTypeInventory && orderType != entities.OrderTypeCustom {
		log.Printf("🏭 [SKIP] Product creation only for INVENTORY and CUSTOM orders, got: %s", orderType)
		return
	}

	// Obtener items de la orden
	items, ok := event.Data["items"].([]entities.OrderItem)
	if !ok {
		log.Printf("🏭 [ERROR] Invalid items in event data")
		return
	}

	log.Printf("🏭 [PRODUCT CREATION] Processing %d items for %s order #%d", len(items), orderType, event.OrderID)

	// Crear productos para cada item
	for _, item := range items {
		if err := h.createProductForItem(ctx, &item, event.OrderID, orderType); err != nil {
			log.Printf("🏭 [ERROR] Failed to create product for item #%d: %v", item.ID, err)
			continue
		}
	}

	log.Printf("🏭 [PRODUCT CREATION] Completed for order #%d", event.OrderID)
}

// createProductForItem crea un producto/variante para un OrderItem o actualiza stock si ya existe
func (h *ProductCreationHandler) createProductForItem(ctx context.Context, item *entities.OrderItem, orderID uint, orderType entities.OrderType) error {
	if item.IsNewVariant() {
		return h.createProductAndVariant(ctx, item, orderID, orderType)
	}

	return h.updateExistingVariantStock(ctx, item, orderType)
}

// updateExistingVariantStock actualiza el stock de una variante existente
func (h *ProductCreationHandler) updateExistingVariantStock(ctx context.Context, item *entities.OrderItem, orderType entities.OrderType) error {
	variant, err := h.productVariantRepo.GetByID(ctx, item.ProductVariantID)
	if err != nil {
		log.Printf("⚠️  [WARNING] Variant #%d not found for OrderItem #%d: %v", item.ProductVariantID, item.ID, err)
		return nil // No bloqueamos el proceso
	}

	quantityToManufacture := item.GetQuantityToManufacture(variant.ReservedStock)

	// Cargar ProductVariant en el item para IsFullyCoveredByStock
	item.ProductVariant = variant

	if item.IsFullyCoveredByStock() {
		log.Printf("✅ [SKIP] Variant #%d: %s | All %d units covered by reserved stock",
			variant.ID, variant.GetFullName(), variant.ReservedStock)
		return nil
	}

	log.Printf("🏭 [MANUFACTURING] Variant #%d: %s | Requested: %d | Reserved: %d | To Manufacture: %d",
		variant.ID, variant.GetFullName(), item.Quantity, variant.ReservedStock, quantityToManufacture)

	// Incrementar stock con lo que se fabricó
	if err := h.productVariantRepo.UpdateStock(ctx, variant.ID, quantityToManufacture); err != nil {
		return err
	}

	// 🔒 RESERVAR lo fabricado SOLO para órdenes CUSTOM (exclusivo para esa orden)
	if orderType == entities.OrderTypeCustom {
		if err := h.productVariantRepo.ReserveStock(ctx, variant.ID, quantityToManufacture); err != nil {
			log.Printf("❌ [ERROR] Failed to reserve manufactured stock for variant #%d: %v", variant.ID, err)
			return err
		}
		log.Printf("✅ [UPDATED] Variant #%d: %s | Stock increased by %d and reserved for CUSTOM order",
			variant.ID, variant.GetFullName(), quantityToManufacture)
	} else {
		log.Printf("✅ [UPDATED] Variant #%d: %s | Stock increased by %d (available for sale)",
			variant.ID, variant.GetFullName(), quantityToManufacture)
	}

	return nil
}

// createProductAndVariant crea un producto base y/o variante según sea necesario
// Casos:
// 1. Producto nuevo + Variante nueva → Crear ambos
// 2. Producto existente + Variante nueva → Crear solo variante
func (h *ProductCreationHandler) createProductAndVariant(ctx context.Context, item *entities.OrderItem, orderID uint, orderType entities.OrderType) error {
	// Validar que tengamos el nombre del producto
	if item.ProductName == "" {
		log.Printf("❌ [ERROR] OrderItem #%d has no ProductName, cannot create product/variant", item.ID)
		return nil // No bloqueamos el proceso
	}

	log.Printf("🏭 [CREATING] Product: %s | Color: %s | Size: %v | Quantity: %d",
		item.ProductName, item.Color, item.SizeID, item.Quantity)

	// 1. Buscar o crear producto base
	product, err := h.findOrCreateProductBase(ctx, item)
	if err != nil {
		return err
	}

	// 2. Crear variante
	reservedStock := 0
	if orderType == entities.OrderTypeCustom {
		reservedStock = item.Quantity // 🔒 Para CUSTOM, reservar todo
	}

	variant := &entities.ProductVariant{
		ProductID:     product.ID,
		Color:         item.Color,
		SizeID:        item.SizeID,
		Stock:         item.Quantity, // Stock inicial = toda la cantidad fabricada
		ReservedStock: reservedStock, // Reservado solo para CUSTOM
		UnitPrice:     item.UnitPrice,
		IsActive:      true,
	}

	if err := h.productVariantRepo.Create(ctx, variant); err != nil {
		return err
	}

	if orderType == entities.OrderTypeCustom {
		log.Printf("✅ [CREATED] Variant #%d: %s with stock %d (all reserved for CUSTOM order)",
			variant.ID, variant.GetFullName(), variant.Stock)
	} else {
		log.Printf("✅ [CREATED] Variant #%d: %s with stock %d (available for sale)",
			variant.ID, variant.GetFullName(), variant.Stock)
	}

	// 3. Actualizar el OrderItem con el ProductVariantID
	item.ProductVariantID = variant.ID
	if err := h.orderItemRepo.Update(ctx, item); err != nil {
		log.Printf("⚠️  [WARNING] Failed to update OrderItem #%d with ProductVariantID: %v", item.ID, err)
		// No retornamos error porque la variante ya fue creada
	}

	log.Printf("🔗 [LINKED] OrderItem #%d → Variant #%d", item.ID, variant.ID)

	return nil
}

// findOrCreateProductBase busca un producto base por nombre o lo crea si no existe
func (h *ProductCreationHandler) findOrCreateProductBase(ctx context.Context, item *entities.OrderItem) (*entities.Product, error) {
	// Buscar producto base por nombre
	products, err := h.productRepo.List(ctx, map[string]interface{}{
		"name": item.ProductName,
	})

	// Si encontramos el producto, retornarlo
	if err == nil && len(products) > 0 {
		log.Printf("📦 [FOUND] Product base: %s (ID: %d)", products[0].Name, products[0].ID)
		return &products[0], nil
	}

	// Si no existe, crear producto base
	log.Printf("🆕 [CREATING] New product base: %s (Category: %d)", item.ProductName, item.CategoryID)

	product := &entities.Product{
		Name:            item.ProductName,
		Description:     "",
		CategoryID:      item.CategoryID, // Usar la categoría del OrderItem
		MaterialCost:    0,
		LaborCost:       0,
		ProductionCost:  0,
		UnitPrice:       item.UnitPrice,
		WholesalePrice:  item.UnitPrice * 0.8,
		MinWholesaleQty: 10,
		MinStock:        5,
		IsActive:        true,
	}

	if err := h.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	log.Printf("✅ [CREATED] Product base #%d: %s", product.ID, product.Name)

	return product, nil
}
