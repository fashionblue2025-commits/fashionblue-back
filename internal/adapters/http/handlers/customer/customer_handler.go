package customer

import (
	"strconv"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/dto"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/customer"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

type CustomerHandler struct {
	createCustomerUC      *customer.CreateCustomerUseCase
	getCustomerUC         *customer.GetCustomerUseCase
	listCustomersUC       *customer.ListCustomersUseCase
	updateCustomerUC      *customer.UpdateCustomerUseCase
	deleteCustomerUC      *customer.DeleteCustomerUseCase
	getHistoryUC          *customer.GetCustomerHistoryUseCase
	createPaymentUC       *customer.CreatePaymentUseCase
	getUpcomingPaymentsUC *customer.GetUpcomingPaymentsUseCase
	getBalanceUC          *customer.GetCustomerBalanceUseCase
	addTransactionUC      *customer.AddTransactionUseCase
}

func NewCustomerHandler(
	createCustomerUC *customer.CreateCustomerUseCase,
	getCustomerUC *customer.GetCustomerUseCase,
	listCustomersUC *customer.ListCustomersUseCase,
	updateCustomerUC *customer.UpdateCustomerUseCase,
	deleteCustomerUC *customer.DeleteCustomerUseCase,
	getHistoryUC *customer.GetCustomerHistoryUseCase,
	createPaymentUC *customer.CreatePaymentUseCase,
	getUpcomingPaymentsUC *customer.GetUpcomingPaymentsUseCase,
	getBalanceUC *customer.GetCustomerBalanceUseCase,
	addTransactionUC *customer.AddTransactionUseCase,
) *CustomerHandler {
	return &CustomerHandler{
		createCustomerUC:      createCustomerUC,
		getCustomerUC:         getCustomerUC,
		listCustomersUC:       listCustomersUC,
		updateCustomerUC:      updateCustomerUC,
		deleteCustomerUC:      deleteCustomerUC,
		getHistoryUC:          getHistoryUC,
		createPaymentUC:       createPaymentUC,
		getUpcomingPaymentsUC: getUpcomingPaymentsUC,
		getBalanceUC:          getBalanceUC,
		addTransactionUC:      addTransactionUC,
	}
}

type CreateCustomerRequest struct {
	Name        string             `json:"name" validate:"required"`
	Phone       string             `json:"phone" validate:"required"`
	Address     string             `json:"address"`
	RiskLevel   entities.RiskLevel `json:"risk_level" validate:"required"`
	ShirtSizeID *uint              `json:"shirt_size_id"`
	PantsSizeID *uint              `json:"pants_size_id"`
	ShoesSizeID *uint              `json:"shoes_size_id"`
	Birthday    *time.Time         `json:"birthday"`
	Notes       string             `json:"notes"`
}

type UpdateCustomerRequest struct {
	Name        string             `json:"name"`
	Phone       string             `json:"phone"`
	Address     string             `json:"address"`
	RiskLevel   entities.RiskLevel `json:"risk_level"`
	ShirtSizeID *uint              `json:"shirt_size_id"`
	PantsSizeID *uint              `json:"pants_size_id"`
	ShoesSizeID *uint              `json:"shoes_size_id"`
	Birthday    *time.Time         `json:"birthday"`
	Notes       string             `json:"notes"`
	IsActive    *bool              `json:"is_active"`
}

func (h *CustomerHandler) Create(c echo.Context) error {
	var req CreateCustomerRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	cust := &entities.Customer{
		Name:        req.Name,
		Phone:       req.Phone,
		Address:     req.Address,
		RiskLevel:   req.RiskLevel,
		ShirtSizeID: req.ShirtSizeID,
		PantsSizeID: req.PantsSizeID,
		ShoesSizeID: req.ShoesSizeID,
		Birthday:    req.Birthday,
		Notes:       req.Notes,
		IsActive:    true,
	}

	if err := h.createCustomerUC.Execute(c.Request().Context(), cust); err != nil {
		return response.BadRequest(c, "Failed to create customer", err)
	}

	// Convertir a DTO
	customerDTO := dto.ToCustomerDTO(cust)

	return response.Created(c, "Customer created successfully", customerDTO)
}

func (h *CustomerHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID", err)
	}

	cust, err := h.getCustomerUC.Execute(c.Request().Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "Customer not found")
	}

	// Convertir a DTO
	customerDTO := dto.ToCustomerDTO(cust)

	return response.OK(c, "Customer retrieved successfully", customerDTO)
}

func (h *CustomerHandler) List(c echo.Context) error {
	filters := make(map[string]interface{})

	if name := c.QueryParam("name"); name != "" {
		filters["name"] = name
	}

	customers, err := h.listCustomersUC.Execute(c.Request().Context(), filters)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve customers", err)
	}

	// Convertir a DTOs
	customerDTOs := dto.ToCustomerDTOList(customers)

	return response.OK(c, "Customers retrieved successfully", customerDTOs)
}

func (h *CustomerHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID", err)
	}

	// Obtener cliente existente
	existingCustomer, err := h.getCustomerUC.Execute(c.Request().Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "Customer not found")
	}

	var req UpdateCustomerRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	// Actualizar solo los campos proporcionados
	if req.Name != "" {
		existingCustomer.Name = req.Name
	}
	if req.Phone != "" {
		existingCustomer.Phone = req.Phone
	}
	if req.Address != "" {
		existingCustomer.Address = req.Address
	}
	if req.RiskLevel != "" {
		existingCustomer.RiskLevel = req.RiskLevel
	}
	if req.ShirtSizeID != nil {
		existingCustomer.ShirtSizeID = req.ShirtSizeID
	}
	if req.PantsSizeID != nil {
		existingCustomer.PantsSizeID = req.PantsSizeID
	}
	if req.ShoesSizeID != nil {
		existingCustomer.ShoesSizeID = req.ShoesSizeID
	}
	if req.Birthday != nil {
		existingCustomer.Birthday = req.Birthday
	}
	if req.Notes != "" {
		existingCustomer.Notes = req.Notes
	}
	if req.IsActive != nil {
		existingCustomer.IsActive = *req.IsActive
	}

	existingCustomer.ID = uint(id)
	if err := h.updateCustomerUC.Execute(c.Request().Context(), existingCustomer); err != nil {
		return response.BadRequest(c, "Failed to update customer", err)
	}

	return response.OK(c, "Customer updated successfully", existingCustomer)
}

func (h *CustomerHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID", err)
	}

	if err := h.deleteCustomerUC.Execute(c.Request().Context(), uint(id)); err != nil {
		return response.BadRequest(c, "Failed to delete customer", err)
	}

	return response.OK(c, "Customer deleted successfully", nil)
}

func (h *CustomerHandler) GetHistory(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID", err)
	}

	history, err := h.getHistoryUC.Execute(c.Request().Context(), uint(id))
	if err != nil {
		return response.InternalServerError(c, "Failed to get customer history", err)
	}

	// Convertir a DTOs
	historyDTOs := dto.ToCustomerTransactionDTOListFromSlice(history)

	return response.OK(c, "Customer history retrieved successfully", historyDTOs)
}

// CreatePayment crea un abono/pago de un cliente
func (h *CustomerHandler) CreatePayment(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID", err)
	}

	type CreatePaymentRequest struct {
		Amount          float64    `json:"amount" validate:"required,gt=0"`
		PaymentMethodID uint       `json:"payment_method_id" validate:"required"`
		Concept         string     `json:"concept" validate:"required"`
		Date            *time.Time `json:"date"`
	}

	var req CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	// Si no se proporciona fecha, usar la actual
	paymentDate := time.Now()
	if req.Date != nil {
		paymentDate = *req.Date
	}

	input := customer.CreatePaymentInput{
		CustomerID:      uint(id),
		Amount:          req.Amount,
		PaymentMethodID: req.PaymentMethodID,
		Concept:         req.Concept,
		Date:            paymentDate,
	}

	transaction, err := h.createPaymentUC.Execute(c.Request().Context(), input)
	if err != nil {
		return response.InternalServerError(c, "Failed to create payment", err)
	}

	// Convertir a DTO
	transactionDTO := dto.ToCustomerTransactionDTO(transaction)

	return response.Created(c, "Payment created successfully", transactionDTO)
}

// GetUpcomingPayments obtiene clientes con pagos próximos
func (h *CustomerHandler) GetUpcomingPayments(c echo.Context) error {
	daysRange := 3 // Por defecto 3 días
	if days := c.QueryParam("days"); days != "" {
		if parsed, err := strconv.Atoi(days); err == nil && parsed > 0 {
			daysRange = parsed
		}
	}

	customers, err := h.getUpcomingPaymentsUC.Execute(c.Request().Context(), daysRange)
	if err != nil {
		return response.InternalServerError(c, "Failed to get upcoming payments", err)
	}

	// Convertir a DTOs
	result := make([]dto.CustomerWithBalanceDTO, len(customers))
	for i, c := range customers {
		result[i] = dto.ToCustomerWithBalanceDTO(&c.Customer, c.Balance)
	}

	return response.OK(c, "Upcoming payments retrieved successfully", result)
}

// GetBalance obtiene el balance actual de un cliente
func (h *CustomerHandler) GetBalance(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID", err)
	}

	balance, err := h.getBalanceUC.Execute(c.Request().Context(), uint(id))
	if err != nil {
		return response.InternalServerError(c, "Failed to get customer balance", err)
	}

	// Convertir a DTO
	balanceDTO := dto.CustomerBalanceDTO{
		CustomerID: uint(id),
		Balance:    balance,
	}

	return response.OK(c, "Balance retrieved successfully", balanceDTO)
}
