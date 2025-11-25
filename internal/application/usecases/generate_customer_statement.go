package usecases

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"github.com/jung-kurt/gofpdf"
)

type GenerateCustomerStatementUseCase struct {
	customerRepo    ports.CustomerRepository
	transactionRepo ports.CustomerTransactionRepository
}

func NewGenerateCustomerStatementUseCase(
	customerRepo ports.CustomerRepository,
	transactionRepo ports.CustomerTransactionRepository,
) *GenerateCustomerStatementUseCase {
	return &GenerateCustomerStatementUseCase{
		customerRepo:    customerRepo,
		transactionRepo: transactionRepo,
	}
}

// StatementRequest representa los parámetros para generar el estado de cuenta
type StatementRequest struct {
	CustomerID uint
	Days       *int // nil = todas las transacciones, valor = últimos X días
}

// StatementResponse contiene el PDF generado
type StatementResponse struct {
	PDFBytes []byte
	Filename string
}

// Execute genera el PDF del estado de cuenta del cliente
func (uc *GenerateCustomerStatementUseCase) Execute(ctx context.Context, req StatementRequest) (*StatementResponse, error) {
	// Obtener el cliente
	customer, err := uc.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo cliente: %w", err)
	}

	// Obtener todas las transacciones del cliente
	allTransactions, err := uc.transactionRepo.ListByCustomer(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo transacciones: %w", err)
	}

	var transactions []entities.CustomerTransaction
	var previousBalance float64
	var cutoffDate time.Time

	// Si se especificaron días, filtrar transacciones
	if req.Days != nil && *req.Days > 0 {
		cutoffDate = time.Now().AddDate(0, 0, -*req.Days)

		// Separar transacciones en "anteriores" y "del período"
		for _, tx := range allTransactions {
			if tx.Date.Before(cutoffDate) {
				// Calcular balance anterior
				if tx.Type == entities.TransactionTypeDebt {
					previousBalance += tx.Amount
				} else {
					previousBalance -= tx.Amount
				}
			} else {
				// Transacción del período a mostrar
				transactions = append(transactions, tx)
			}
		}
	} else {
		// Todas las transacciones
		transactions = allTransactions
		previousBalance = 0
	}

	// Ordenar transacciones de forma ascendente por fecha (más antiguas primero)
	// Esto es solo para el PDF, no afecta otras vistas
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Date.Before(transactions[j].Date)
	})

	// Generar el PDF con soporte UTF-8
	pdf := gofpdf.New("P", "mm", "Letter", "")

	// Configurar para usar UTF-8
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.AddPage()

	// Configurar fuentes y color del título (Azul Principal)
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(29, 161, 242) // #1DA1F2 - Azul Principal

	// Título
	pdf.Cell(0, 10, tr("ESTADO DE CUENTA"))
	pdf.Ln(12)

	// Resetear color a negro para el resto
	pdf.SetTextColor(0, 0, 0)

	// Información del cliente
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 6, tr(fmt.Sprintf("Cliente: %s", customer.Name)))
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	if customer.Phone != "" {
		pdf.Cell(0, 5, tr(fmt.Sprintf("Teléfono: %s", customer.Phone)))
		pdf.Ln(5)
	}
	if customer.Address != "" {
		pdf.Cell(0, 5, tr(fmt.Sprintf("Dirección: %s", customer.Address)))
		pdf.Ln(5)
	}

	// Fecha del reporte
	pdf.SetFont("Arial", "", 9)
	pdf.Cell(0, 5, tr(fmt.Sprintf("Fecha de generación: %s", time.Now().Format("02/01/2006 15:04"))))
	pdf.Ln(5)

	if req.Days != nil && *req.Days > 0 {
		pdf.Cell(0, 5, tr(fmt.Sprintf("Período: Últimos %d días (desde %s)", *req.Days, cutoffDate.Format("02/01/2006"))))
		pdf.Ln(5)
	} else {
		pdf.Cell(0, 5, tr("Período: Todas las transacciones"))
		pdf.Ln(5)
	}

	pdf.Ln(5)

	// Encabezado de la tabla con colores Fashion Blue
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(29, 161, 242)  // #1DA1F2 - Azul Principal
	pdf.SetTextColor(255, 255, 255) // Texto blanco

	pdf.CellFormat(25, 7, tr("Fecha"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 7, tr("Tipo"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(80, 7, tr("Descripción"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 7, tr("Débito"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 7, tr("Crédito"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 7, tr("Saldo"), "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	// Resetear color de texto a negro
	pdf.SetTextColor(0, 0, 0)

	// Balance anterior (si aplica)
	currentBalance := previousBalance
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(243, 229, 245) // #F3E5F5 - Púrpura muy claro

	if previousBalance != 0 {
		pdf.SetTextColor(156, 39, 176) // #9C27B0 - Púrpura (para Balance Anterior)
		pdf.CellFormat(25, 6, "", "1", 0, "L", true, 0, "")
		pdf.CellFormat(20, 6, "", "1", 0, "L", true, 0, "")
		pdf.CellFormat(80, 6, tr("Balance Anterior"), "1", 0, "L", true, 0, "")
		pdf.CellFormat(25, 6, "", "1", 0, "R", true, 0, "")
		pdf.CellFormat(25, 6, "", "1", 0, "R", true, 0, "")
		pdf.CellFormat(25, 6, formatCurrency(previousBalance), "1", 0, "R", true, 0, "")
		pdf.Ln(-1)
		pdf.SetTextColor(0, 0, 0) // Resetear a negro
	}

	// Transacciones
	pdf.SetFont("Arial", "", 8)

	for _, tx := range transactions {
		// Calcular saldo acumulado
		if tx.Type == entities.TransactionTypeDebt {
			currentBalance += tx.Amount
		} else {
			currentBalance -= tx.Amount
		}

		// Fecha
		pdf.CellFormat(25, 6, tx.Date.Format("02/01/2006"), "1", 0, "C", false, 0, "")

		// Tipo
		typeText := "DEUDA"
		if tx.Type == entities.TransactionTypePayment {
			typeText = "ABONO"
		}
		pdf.CellFormat(20, 6, typeText, "1", 0, "C", false, 0, "")

		// Descripción (truncar si es muy larga)
		description := tx.Description
		if len(description) == 0 {
			description = "-"
		}
		if len(description) > 50 {
			description = description[:47] + "..."
		}
		pdf.CellFormat(80, 6, tr(description), "1", 0, "L", false, 0, "")

		// Débito (DEUDA)
		debitText := ""
		if tx.Type == entities.TransactionTypeDebt {
			debitText = formatCurrency(tx.Amount)
		}
		pdf.CellFormat(25, 6, debitText, "1", 0, "R", false, 0, "")

		// Crédito (ABONO)
		creditText := ""
		if tx.Type == entities.TransactionTypePayment {
			creditText = formatCurrency(tx.Amount)
		}
		pdf.CellFormat(25, 6, creditText, "1", 0, "R", false, 0, "")

		// Saldo
		pdf.CellFormat(25, 6, formatCurrency(currentBalance), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}

	// Total final con color secundario (Rosa/Magenta)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(233, 30, 99)   // #E91E63 - Rosa/Magenta
	pdf.SetTextColor(255, 255, 255) // Texto blanco
	pdf.CellFormat(150, 8, tr("SALDO ACTUAL"), "1", 0, "R", true, 0, "")
	pdf.CellFormat(50, 8, formatCurrency(currentBalance), "1", 0, "R", true, 0, "")
	pdf.Ln(-1)

	// Resetear color de texto a negro
	pdf.SetTextColor(0, 0, 0)

	// Resumen
	pdf.Ln(5)
	pdf.SetFont("Arial", "", 9)

	totalDebitos := previousBalance
	totalCreditos := 0.0

	for _, tx := range transactions {
		if tx.Type == entities.TransactionTypeDebt {
			totalDebitos += tx.Amount
		} else {
			totalCreditos += tx.Amount
		}
	}

	if req.Days != nil && *req.Days > 0 {
		pdf.Cell(0, 5, tr(fmt.Sprintf("Balance anterior: %s", formatCurrency(previousBalance))))
		pdf.Ln(5)
		pdf.Cell(0, 5, tr(fmt.Sprintf("Total débitos en el período: %s", formatCurrency(totalDebitos-previousBalance))))
		pdf.Ln(5)
		pdf.Cell(0, 5, tr(fmt.Sprintf("Total créditos en el período: %s", formatCurrency(totalCreditos))))
		pdf.Ln(5)
	} else {
		pdf.Cell(0, 5, tr(fmt.Sprintf("Total débitos: %s", formatCurrency(totalDebitos))))
		pdf.Ln(5)
		pdf.Cell(0, 5, tr(fmt.Sprintf("Total créditos: %s", formatCurrency(totalCreditos))))
		pdf.Ln(5)
	}

	// Generar bytes del PDF
	writer := &bytes.Buffer{}
	err = pdf.Output(writer)
	if err != nil {
		return nil, fmt.Errorf("error generando PDF: %w", err)
	}

	// Generar nombre de archivo
	filename := fmt.Sprintf("estado_cuenta_%s_%s.pdf",
		sanitizeFilename(customer.Name),
		time.Now().Format("20060102"))

	return &StatementResponse{
		PDFBytes: writer.Bytes(),
		Filename: filename,
	}, nil
}

// formatCurrency formatea un número como moneda
func formatCurrency(amount float64) string {
	// Formatear con separador de miles
	strAmount := fmt.Sprintf("%.0f", amount)

	// Agregar separadores de miles
	if len(strAmount) <= 3 {
		return "$" + strAmount
	}

	result := ""
	for i, digit := range strAmount {
		if i > 0 && (len(strAmount)-i)%3 == 0 {
			result += "."
		}
		result += string(digit)
	}

	return "$" + result
}

// sanitizeFilename limpia el nombre del archivo
func sanitizeFilename(s string) string {
	// Reemplazar espacios y caracteres especiales
	result := ""
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			result += string(r)
		} else if r == ' ' {
			result += "_"
		}
	}
	return result
}
