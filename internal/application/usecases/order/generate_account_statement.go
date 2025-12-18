package order

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

type GenerateAccountStatementUseCase struct {
	orderRepository ports.OrderRepository
}

func NewGenerateAccountStatementUseCase(orderRepository ports.OrderRepository) *GenerateAccountStatementUseCase {
	return &GenerateAccountStatementUseCase{
		orderRepository: orderRepository,
	}
}

type AccountStatementData struct {
	OrderID         uint
	OrderNumber     string
	StatementNumber string
	SellerName      string
	SellerID        string
	ClientName      string
	City            string
	Date            time.Time
	Concept         string
	TotalAmount     float64
	BankAccount     string
}

// GetDraft obtiene los datos del borrador de cuenta de cobro para que el usuario los edite
func (uc *GenerateAccountStatementUseCase) GetDraft(ctx context.Context, orderID uint) (*AccountStatementData, error) {
	// Obtener la orden
	order, err := uc.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Generar número de cuenta de cobro basado en el order number
	statementNumber := fmt.Sprintf("%s-CC", order.OrderNumber)

	// Construir concepto basado en los items de la orden
	concept := uc.buildConcept(order)

	data := AccountStatementData{
		OrderID:         orderID,
		OrderNumber:     order.OrderNumber,
		StatementNumber: statementNumber,
		SellerName:      "SONIA PATRICIA ORTIZ",
		SellerID:        "30323685",
		ClientName:      order.CustomerName,
		City:            "MANIZALES",
		Date:            time.Now(),
		Concept:         concept,
		TotalAmount:     order.TotalAmount,
		BankAccount:     "3122684372", // Este valor podría venir de configuración
	}

	return &data, nil
}

// GeneratePDF genera el PDF de la cuenta de cobro con los datos confirmados
func (uc *GenerateAccountStatementUseCase) GeneratePDF(ctx context.Context, data AccountStatementData) ([]byte, error) {
	return uc.generatePDF(data)
}

func (uc *GenerateAccountStatementUseCase) buildConcept(order *entities.Order) string {
	// Construir concepto basado en los items
	if len(order.Items) == 0 {
		return fmt.Sprintf("Venta de productos según orden %s", order.OrderNumber)
	}

	// Tomar los primeros items para el concepto
	itemsCount := len(order.Items)
	if itemsCount == 1 {
		return fmt.Sprintf("Venta de %d chaquetas, entregadas a la %s, conforme a acuerdo comercial.",
			order.Items[0].Quantity, order.CustomerName)
	}

	totalQuantity := 0
	for _, item := range order.Items {
		totalQuantity += item.Quantity
	}

	return fmt.Sprintf("Venta de %d unidades de productos diversos, entregadas según orden %s, conforme a acuerdo comercial.",
		totalQuantity, order.OrderNumber)
}

func (uc *GenerateAccountStatementUseCase) generatePDF(data AccountStatementData) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()
	pdf.SetMargins(20, 20, 20)

	// Configurar traductor para caracteres especiales (tildes)
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// Formatear fecha en español
	dateFormatted := uc.formatDateSpanish(data.Date)

	// Título
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, tr(fmt.Sprintf("CUENTA DE COBRO N° %s", data.StatementNumber)))
	pdf.Ln(15)

	// Información del vendedor
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 6, tr(data.SellerName))
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, tr(fmt.Sprintf("CLIENTE: %s", data.ClientName)))
	pdf.Ln(6)
	pdf.Cell(0, 6, tr(fmt.Sprintf("CIUDAD Y FECHA: %s, %s", data.City, dateFormatted)))
	pdf.Ln(12)

	// Texto legal
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, tr("PARA EFECTOS LEGALES RELACIONADOS CON LAS NORMAS FISCALES ME"))
	pdf.Ln(6)
	pdf.Cell(0, 6, tr("PERMITO INFORMAR:"))
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	legalText := []string{
		"1. Soy residente en Colombia: SI",
		"2. Pertenezco al régimen no responsable de IVA.",
		"3. Número de trabajadores a mi cargo o contratistas que tengo vinculados a mi",
		"   actividad es: (0)",
		"4. Que me encuentro dentro de las situaciones contempladas en el Art del Decreto",
		"   1165 de 1996 (Obligados a Facturar)",
	}

	for _, line := range legalText {
		pdf.Cell(0, 5, tr(line))
		pdf.Ln(5)
	}

	pdf.Ln(5)
	pdf.MultiCell(0, 5, tr("En consecuencia y de conformidad con la normatividad fiscal, no estoy obligado a expedir factura ni documento equivalente."), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 5, tr("Por lo anterior, para tener en cuenta en el proceso de retención en la fuente para empleados según lo establecido en la ley 1819 de 2016."), "", "", false)
	pdf.Ln(10)

	// Concepto
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, tr("POR CONCEPTO DE:"))
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, tr(data.Concept), "", "", false)
	pdf.Ln(10)

	// Valor total con formato colombiano
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, tr(fmt.Sprintf("VALOR TOTAL %s", uc.formatCOP(data.TotalAmount))))
	pdf.Ln(12)

	// Información bancaria
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, tr(fmt.Sprintf("Favor consignar a mi cuenta de nequi %s", data.BankAccount)))
	pdf.Ln(6)
	pdf.Cell(0, 6, tr("Atentamente,"))
	pdf.Ln(15)

	// Firma
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, tr(data.SellerName))
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, data.SellerID)

	// Generar buffer con el PDF
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// formatDateSpanish formatea una fecha en español (ej: "15 DE ENERO DE 2025")
func (uc *GenerateAccountStatementUseCase) formatDateSpanish(date time.Time) string {
	months := []string{"", "ENERO", "FEBRERO", "MARZO", "ABRIL", "MAYO", "JUNIO", "JULIO", "AGOSTO", "SEPTIEMBRE", "OCTUBRE", "NOVIEMBRE", "DICIEMBRE"}
	return fmt.Sprintf("%d DE %s DE %d", date.Day(), months[date.Month()], date.Year())
}

// formatCOP formatea un número como pesos colombianos ($1.000.000)
func (uc *GenerateAccountStatementUseCase) formatCOP(amount float64) string {
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
