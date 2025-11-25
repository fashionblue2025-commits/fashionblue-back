package financial_transaction

import (
	"strconv"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/dto"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/financial_transaction"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

type FinancialTransactionHandler struct {
	createTransactionUC *financial_transaction.CreateTransactionUseCase
	getTransactionUC    *financial_transaction.GetTransactionUseCase
	listTransactionsUC  *financial_transaction.ListTransactionsUseCase
	getBalanceUC        *financial_transaction.GetBalanceUseCase
}

func NewFinancialTransactionHandler(
	createTransactionUC *financial_transaction.CreateTransactionUseCase,
	getTransactionUC *financial_transaction.GetTransactionUseCase,
	listTransactionsUC *financial_transaction.ListTransactionsUseCase,
	getBalanceUC *financial_transaction.GetBalanceUseCase,
) *FinancialTransactionHandler {
	return &FinancialTransactionHandler{
		createTransactionUC: createTransactionUC,
		getTransactionUC:    getTransactionUC,
		listTransactionsUC:  listTransactionsUC,
		getBalanceUC:        getBalanceUC,
	}
}

func (h *FinancialTransactionHandler) Create(c echo.Context) error {
	var transaction entities.FinancialTransaction
	if err := c.Bind(&transaction); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := h.createTransactionUC.Execute(c.Request().Context(), &transaction); err != nil {
		return response.BadRequest(c, "Failed to create transaction", err)
	}

	return response.Created(c, "Transaction created successfully", dto.ToFinancialTransactionDTO(&transaction))
}

func (h *FinancialTransactionHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid transaction ID", err)
	}

	transaction, err := h.getTransactionUC.Execute(c.Request().Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "Transaction not found")
	}

	return response.OK(c, "Transaction retrieved successfully", dto.ToFinancialTransactionDTO(transaction))
}

func (h *FinancialTransactionHandler) List(c echo.Context) error {
	filters := make(map[string]interface{})

	// Filtros opcionales
	if transactionType := c.QueryParam("type"); transactionType != "" {
		filters["type"] = transactionType
	}
	if category := c.QueryParam("category"); category != "" {
		filters["category"] = category
	}
	if startDate := c.QueryParam("start_date"); startDate != "" {
		filters["start_date"] = startDate
	}
	if endDate := c.QueryParam("end_date"); endDate != "" {
		filters["end_date"] = endDate
	}

	transactions, err := h.listTransactionsUC.Execute(c.Request().Context(), filters)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve transactions", err)
	}

	return response.OK(c, "Transactions retrieved successfully", dto.ToFinancialTransactionDTOList(transactions))
}

func (h *FinancialTransactionHandler) GetBalance(c echo.Context) error {
	balance, err := h.getBalanceUC.Execute(c.Request().Context())
	if err != nil {
		return response.InternalServerError(c, "Failed to get balance", err)
	}

	return response.OK(c, "Balance retrieved successfully", balance)
}
