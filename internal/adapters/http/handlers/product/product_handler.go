package product

import (
	"io"
	"strconv"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/dto"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/product"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	createProductUC        *product.CreateProductUseCase
	getProductUC           *product.GetProductUseCase
	listProductsUC         *product.ListProductsUseCase
	updateProductUC        *product.UpdateProductUseCase
	deleteProductUC        *product.DeleteProductUseCase
	getLowStockUC          *product.GetLowStockProductsUseCase
	uploadPhotoUC          *product.UploadProductPhotoUseCase
	uploadMultiplePhotosUC *product.UploadMultiplePhotosUseCase
	getPhotosUC            *product.GetProductPhotosUseCase
	deletePhotoUC          *product.DeleteProductPhotoUseCase
	setPrimaryPhotoUC      *product.SetPrimaryPhotoUseCase
}

func NewProductHandler(
	createProductUC *product.CreateProductUseCase,
	getProductUC *product.GetProductUseCase,
	listProductsUC *product.ListProductsUseCase,
	updateProductUC *product.UpdateProductUseCase,
	deleteProductUC *product.DeleteProductUseCase,
	getLowStockUC *product.GetLowStockProductsUseCase,
	uploadPhotoUC *product.UploadProductPhotoUseCase,
	uploadMultiplePhotosUC *product.UploadMultiplePhotosUseCase,
	getPhotosUC *product.GetProductPhotosUseCase,
	deletePhotoUC *product.DeleteProductPhotoUseCase,
	setPrimaryPhotoUC *product.SetPrimaryPhotoUseCase,
) *ProductHandler {
	return &ProductHandler{
		createProductUC:        createProductUC,
		getProductUC:           getProductUC,
		listProductsUC:         listProductsUC,
		updateProductUC:        updateProductUC,
		deleteProductUC:        deleteProductUC,
		getLowStockUC:          getLowStockUC,
		uploadPhotoUC:          uploadPhotoUC,
		uploadMultiplePhotosUC: uploadMultiplePhotosUC,
		getPhotosUC:            getPhotosUC,
		deletePhotoUC:          deletePhotoUC,
		setPrimaryPhotoUC:      setPrimaryPhotoUC,
	}
}

func (h *ProductHandler) Create(c echo.Context) error {
	var req dto.ProductDTO
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	// Convertir DTO a entidad
	product := &entities.Product{
		Name:            req.Name,
		Description:     req.Description,
		CategoryID:      req.CategoryID,
		MaterialCost:    req.MaterialCost,
		LaborCost:       req.LaborCost,
		ProductionCost:  req.ProductionCost,
		UnitPrice:       req.UnitPrice,
		WholesalePrice:  req.WholesalePrice,
		MinWholesaleQty: req.MinWholesaleQty,
		MinStock:        req.MinStock,
		IsActive:        req.IsActive,
	}

	// Convertir variantes si existen
	if len(req.Variants) > 0 {
		product.Variants = make([]entities.ProductVariant, len(req.Variants))
		for i, variantDTO := range req.Variants {
			product.Variants[i] = entities.ProductVariant{
				Color:     variantDTO.Color,
				SizeID:    variantDTO.SizeID,
				Stock:     variantDTO.Stock,
				UnitPrice: variantDTO.UnitPrice,
				IsActive:  variantDTO.IsActive,
			}
		}
	}

	if err := h.createProductUC.Execute(c.Request().Context(), product); err != nil {
		return response.BadRequest(c, "Failed to create product", err)
	}

	return response.Created(c, "Product created successfully", dto.ToProductDTO(product))
}

func (h *ProductHandler) GetByID(c echo.Context) error {

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID", err)
	}

	product, err := h.getProductUC.Execute(c.Request().Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "Product not found")
	}

	return response.OK(c, "Product retrieved successfully", dto.ToProductDTO(product))
}

func (h *ProductHandler) List(c echo.Context) error {
	filters := make(map[string]interface{})

	if categoryID := c.QueryParam("category_id"); categoryID != "" {
		if id, err := strconv.ParseUint(categoryID, 10, 32); err == nil {
			filters["category_id"] = uint(id)
		}
	}

	if name := c.QueryParam("name"); name != "" {
		filters["name"] = name
	}

	products, err := h.listProductsUC.Execute(c.Request().Context(), filters)
	if err != nil {
		return response.InternalServerError(c, "Failed to list products", err)
	}

	return response.OK(c, "Products retrieved successfully", dto.ToProductDTOList(products))
}

func (h *ProductHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID", err)
	}

	var product entities.Product
	if err := c.Bind(&product); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	product.ID = uint(id)
	if err := h.updateProductUC.Execute(c.Request().Context(), &product); err != nil {
		return response.BadRequest(c, "Failed to update product", err)
	}

	return response.OK(c, "Product updated successfully", dto.ToProductDTO(&product))
}

func (h *ProductHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID", err)
	}

	if err := h.deleteProductUC.Execute(c.Request().Context(), uint(id)); err != nil {
		return response.BadRequest(c, "Failed to delete product", err)
	}

	return response.OK(c, "Product deleted successfully", nil)
}

func (h *ProductHandler) GetLowStock(c echo.Context) error {
	products, err := h.getLowStockUC.Execute(c.Request().Context())
	if err != nil {
		return response.InternalServerError(c, "Failed to get low stock products", err)
	}

	return response.OK(c, "Low stock products retrieved successfully", dto.ToProductDTOList(products))
}

// UploadPhoto sube una o múltiples fotos para un producto
func (h *ProductHandler) UploadPhotos(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID", err)
	}

	// Obtener el formulario multipart
	form, err := c.MultipartForm()
	if err != nil {
		return response.BadRequest(c, "Invalid form data", err)
	}

	// Obtener archivos (puede ser uno o múltiples con el mismo nombre "photos")
	files := form.File["files"]
	if len(files) == 0 {
		return response.BadRequest(c, "At least one photo is required", nil)
	}

	// Preparar array de fotos para subir
	photosToUpload := make([]product.PhotoUpload, 0, len(files))

	for i, file := range files {
		// Abrir archivo
		src, err := file.Open()
		if err != nil {
			return response.InternalServerError(c, "Failed to open file", err)
		}

		// Leer contenido
		fileData, err := io.ReadAll(src)
		src.Close()
		if err != nil {
			return response.InternalServerError(c, "Failed to read file", err)
		}

		// Primera foto es primary por defecto
		isPrimary := i == 0

		photosToUpload = append(photosToUpload, product.PhotoUpload{
			FileName:    file.Filename,
			FileData:    fileData,
			ContentType: file.Header.Get("Content-Type"),
			Description: c.FormValue("description"),
			IsPrimary:   isPrimary,
		})
	}

	// Ejecutar caso de uso
	uploadedPhotos, err := h.uploadMultiplePhotosUC.Execute(
		c.Request().Context(),
		uint(productID),
		photosToUpload,
	)
	if err != nil {
		return response.BadRequest(c, "Failed to upload photos", err)
	}

	// Convertir a DTOs
	photoDTOs := dto.ToProductPhotoDTOList(uploadedPhotos)
	return response.Created(c, "Photos uploaded successfully", photoDTOs)
}

// GetPhotos obtiene todas las fotos de un producto
func (h *ProductHandler) GetPhotos(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID", err)
	}

	photos, err := h.getPhotosUC.Execute(c.Request().Context(), uint(productID))
	if err != nil {
		return response.InternalServerError(c, "Failed to get photos", err)
	}

	// Convertir a DTOs
	photoDTOs := dto.ToProductPhotoDTOList(photos)
	return response.OK(c, "Photos retrieved successfully", photoDTOs)
}

// DeletePhoto elimina una foto de un producto
func (h *ProductHandler) DeletePhoto(c echo.Context) error {
	photoID, err := strconv.ParseUint(c.Param("photoId"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid photo ID", err)
	}

	if err := h.deletePhotoUC.Execute(c.Request().Context(), uint(photoID)); err != nil {
		return response.BadRequest(c, "Failed to delete photo", err)
	}

	return response.OK(c, "Photo deleted successfully", nil)
}

// SetPrimaryPhoto establece una foto como principal
func (h *ProductHandler) SetPrimaryPhoto(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID", err)
	}

	photoID, err := strconv.ParseUint(c.Param("photoId"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid photo ID", err)
	}

	if err := h.setPrimaryPhotoUC.Execute(c.Request().Context(), uint(photoID), uint(productID)); err != nil {
		return response.BadRequest(c, "Failed to set primary photo", err)
	}

	return response.OK(c, "Primary photo set successfully", nil)
}

// GetStats eliminado - dependía de Sales que ya no existe
