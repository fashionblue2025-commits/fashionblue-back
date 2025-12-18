package financial_transaction

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"github.com/jung-kurt/gofpdf"
)

type GeneratePDFUseCase struct {
	transactionRepo ports.FinancialTransactionRepository
}

func NewGeneratePDFUseCase(transactionRepo ports.FinancialTransactionRepository) *GeneratePDFUseCase {
	return &GeneratePDFUseCase{
		transactionRepo: transactionRepo,
	}
}

func (uc *GeneratePDFUseCase) Execute(ctx context.Context, filters map[string]interface{}) ([]byte, error) {
	// Obtener transacciones según filtros
	transactions, err := uc.transactionRepo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Crear PDF con soporte UTF-8
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()

	// Configurar traductor para caracteres especiales (tildes)
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// Configurar fuentes
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, tr("Reporte de Transacciones Financieras"))
	pdf.Ln(12)

	// Información de filtros aplicados
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, tr(fmt.Sprintf("Fecha de generación: %s", time.Now().Format("02/01/2006 15:04"))))
	pdf.Ln(6)

	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		pdf.Cell(0, 6, tr(fmt.Sprintf("Desde: %s", startDate)))
		pdf.Ln(6)
	}
	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		pdf.Cell(0, 6, tr(fmt.Sprintf("Hasta: %s", endDate)))
		pdf.Ln(6)
	}
	if transactionType, ok := filters["type"].(string); ok && transactionType != "" {
		pdf.Cell(0, 6, tr(fmt.Sprintf("Tipo: %s", transactionType)))
		pdf.Ln(6)
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		pdf.Cell(0, 6, tr(fmt.Sprintf("Categoría: %s", category)))
		pdf.Ln(6)
	}

	pdf.Ln(8)

	// Tabla de transacciones
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(200, 220, 255)

	// Encabezados
	pdf.CellFormat(25, 7, tr("Fecha"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 7, tr("Tipo"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 7, tr("Categoría"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 7, tr("Monto"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(80, 7, tr("Descripción"), "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	// Datos
	pdf.SetFont("Arial", "", 9)
	var totalIncome, totalExpense float64

	for _, t := range transactions {
		// Guardar posición inicial
		x := pdf.GetX()
		y := pdf.GetY()

		// Calcular altura necesaria para la descripción
		// Ancho de la columna descripción = 80mm
		descriptionHeight := pdf.SplitLines([]byte(tr(t.Description)), 80)
		cellHeight := float64(len(descriptionHeight)) * 5 // 5mm por línea
		if cellHeight < 6 {
			cellHeight = 6 // Altura mínima
		}

		// Tipo de transacción
		typeText := "INGRESO"
		if t.Type == entities.FinancialTransactionTypeExpense {
			typeText = "GASTO"
			totalExpense += t.Amount
		} else {
			totalIncome += t.Amount
		}

		// Fecha
		pdf.Rect(x, y, 25, cellHeight, "D")
		pdf.SetXY(x, y+(cellHeight-6)/2) // Centrar verticalmente
		pdf.Cell(25, 6, t.Date.Format("02/01/2006"))

		// Tipo
		x += 25
		pdf.Rect(x, y, 25, cellHeight, "D")
		pdf.SetXY(x, y+(cellHeight-6)/2) // Centrar verticalmente
		pdf.Cell(25, 6, typeText)

		// Categoría
		x += 25
		pdf.Rect(x, y, 35, cellHeight, "D")
		pdf.SetXY(x, y+(cellHeight-6)/2) // Centrar verticalmente
		pdf.Cell(35, 6, tr(string(t.Category)))

		// Monto
		x += 35
		pdf.Rect(x, y, 35, cellHeight, "D")
		pdf.SetXY(x, y+(cellHeight-6)/2) // Centrar verticalmente
		pdf.CellFormat(35, 6, uc.formatCOP(t.Amount), "", 0, "R", false, 0, "")

		// Descripción con MultiCell para ajuste automático
		x += 35
		pdf.Rect(x, y, 80, cellHeight, "D")
		pdf.SetXY(x, y)
		pdf.MultiCell(80, 5, tr(t.Description), "", "L", false)

		// Mover a la siguiente fila
		pdf.SetXY(10, y+cellHeight)
	}

	// Totales
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 7, tr(fmt.Sprintf("Total Ingresos: %s", uc.formatCOP(totalIncome))))
	pdf.Ln(7)
	pdf.Cell(0, 7, tr(fmt.Sprintf("Total Gastos: %s", uc.formatCOP(totalExpense))))
	pdf.Ln(7)
	pdf.SetFont("Arial", "B", 12)
	balance := totalIncome - totalExpense
	balanceText := tr(fmt.Sprintf("Balance: %s", uc.formatCOP(balance)))
	if balance < 0 {
		pdf.SetTextColor(255, 0, 0)
	} else {
		pdf.SetTextColor(0, 128, 0)
	}
	pdf.Cell(0, 7, balanceText)

	// Generar buffer con el PDF
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// formatCOP formatea un número como pesos colombianos ($1.000.000)
func (uc *GeneratePDFUseCase) formatCOP(amount float64) string {
	// Convertir a entero para evitar decimales
	amountInt := int64(amount)

	// Convertir a string
	amountStr := fmt.Sprintf("%d", amountInt)

	// Si es negativo, guardar el signo
	isNegative := false
	if amountInt < 0 {
		isNegative = true
		amountStr = amountStr[1:] // Quitar el signo menos
	}

	// Agregar separador de miles (punto)
	var result strings.Builder
	length := len(amountStr)

	for i, digit := range amountStr {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteString(".")
		}
		result.WriteRune(digit)
	}

	// Agregar signo de pesos y negativo si aplica
	if isNegative {
		return "-$" + result.String()
	}
	return "$" + result.String()
}
