package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

const (
	// ID del Google Spreadsheet
	SPREADSHEET_ID = "1NfZqCG9bDw4Nmi4nMBHcQbkorTFGY84P7ZqPnMCiCMY"

	// GID de la hoja "Clientes" principal
	// IMPORTANTE: Si no funciona, verifica el GID correcto de tu hoja "Clientes"
	GID_CLIENTES = "2138611282"
)

// CustomerSheetRow representa una fila del Google Sheet principal
type CustomerSheetRow struct {
	Name     string
	SheetGID string // GID extraído del link
}

// CustomerDetailRow representa los datos de la hoja individual del cliente
type CustomerDetailRow struct {
	Name         string
	Email        string
	Phone        string
	Birthday     string
	Transactions []TransactionRow
}

// TransactionRow representa una transacción del historial del cliente
type TransactionRow struct {
	Date          string // Columna A
	Type          string // Columna B (ABONO o DEUDA)
	Size          string // Columna C (Talla, puede estar vacía)
	Concept       string // Columna D (Concepto/Descripción)
	PaymentMethod string // Columna E (Tipo de pago, puede estar vacía)
	Amount        string // Columna F (Valor)
}

func main() {
	fmt.Println("🔄 Iniciando importación de clientes desde Google Sheets...")
	fmt.Println("📊 Spreadsheet ID:", SPREADSHEET_ID)
	fmt.Println("📄 Sheet GID Clientes:", GID_CLIENTES)
	fmt.Println("")

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No se encontró archivo .env, usando variables del sistema")
	}

	// Conectar a la base de datos
	db, err := connectDB()
	if err != nil {
		log.Fatalf("❌ Error conectando a la base de datos: %v", err)
	}

	// Leer lista de clientes desde la hoja principal
	fmt.Println("📥 Paso 1: Descargando lista de clientes...")
	customerList, err := fetchCustomerList()
	if err != nil {
		log.Fatalf("❌ Error obteniendo lista de clientes: %v", err)
	}

	fmt.Printf("📋 Se encontraron %d clientes en la lista\n\n", len(customerList))

	if len(customerList) == 0 {
		fmt.Println("⚠️  No se encontraron clientes para importar.")
		return
	}

	// Importar cada cliente
	imported := 0
	skipped := 0
	errors := 0

	for i, customerRef := range customerList {
		fmt.Printf("\n[%d/%d] Procesando: %s (GID: %s)\n", i+1, len(customerList), customerRef.Name, customerRef.SheetGID)

		// Obtener detalles del cliente desde su hoja individual
		fmt.Printf("   📄 Descargando hoja individual...\n")
		details, err := fetchCustomerDetails(customerRef.SheetGID)
		if err != nil {
			log.Printf("   ❌ Error obteniendo detalles: %v", err)
			errors++
			continue
		}

		// Verificar si el cliente ya existe
		exists, existingID, err := customerExists(db, details.Name, details.Phone)
		if err != nil {
			log.Printf("   ⚠️  Error verificando existencia: %v", err)
			errors++
			continue
		}

		if exists {
			fmt.Printf("   ⏭️  Cliente ya existe (ID: %d), saltando...\n", existingID)
			skipped++
			continue
		}

		// Crear cliente
		if err := createCustomer(db, details); err != nil {
			log.Printf("   ❌ Error creando cliente: %v", err)
			errors++
			continue
		}

		fmt.Printf("   ✅ Cliente creado exitosamente\n")
		imported++
	}

	// Resumen
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📊 RESUMEN DE IMPORTACIÓN")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("✅ Importados: %d\n", imported)
	fmt.Printf("⏭️  Saltados:   %d\n", skipped)
	fmt.Printf("❌ Errores:    %d\n", errors)
	fmt.Printf("📋 Total:      %d\n", len(customerList))
	fmt.Println(strings.Repeat("=", 50))

	if imported > 0 {
		fmt.Println("\n🎉 ¡Importación completada exitosamente!")
	}
}

// fetchCustomerList descarga la hoja "Clientes" y extrae nombres + GIDs
func fetchCustomerList() ([]CustomerSheetRow, error) {
	url := fmt.Sprintf(
		"https://docs.google.com/spreadsheets/d/%s/export?format=csv&gid=%s",
		SPREADSHEET_ID,
		GID_CLIENTES,
	)

	fmt.Printf("📡 Descargando desde: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error descargando sheet: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error HTTP: status code %d - verifica que el sheet sea público", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)

	// Leer todas las filas
	var customers []CustomerSheetRow
	lineNum := 0

	only := []string{"Julián Andrés Trujillo"}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("⚠️  Error en línea %d: %v", lineNum, err)
			lineNum++
			continue
		}

		lineNum++

		// Saltar headers
		if lineNum <= 2 || len(record) < 2 {
			continue
		}

		// Columna A: Nombre
		name := strings.TrimSpace(record[0])
		if name == "" || strings.EqualFold(name, "NOMBRE") {
			continue
		}

		// Columna E (índice 4): FORMULATEXT con el hipervínculo completo
		var gid string
		if len(record) > 4 {
			formulaText := strings.TrimSpace(record[4])
			gid = extractGID(formulaText)
		}

		// Si no está en columna E, intentar columna B (por compatibilidad)
		if gid == "" && len(record) > 1 {
			linkText := strings.TrimSpace(record[1])
			gid = extractGID(linkText)
		}

		// Intentar columna C también
		if gid == "" && len(record) > 2 {
			gid = extractGID(strings.TrimSpace(record[2]))
		}

		if gid == "" {
			log.Printf("  Línea %d: No se pudo extraer GID para '%s'", lineNum, name)
			continue
		}

		if len(only) > 0 && !contains(only, name) {
			continue
		}

		customers = append(customers, CustomerSheetRow{
			Name:     name,
			SheetGID: gid,
		})
	}

	return customers, nil
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// extractGID extrae el GID de una URL de Google Sheets
func extractGID(text string) string {
	// Buscar patrón: gid=NUMEROS
	re := regexp.MustCompile(`gid=(\d+)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// fetchCustomerDetails descarga la hoja individual y extrae los datos del cliente
func fetchCustomerDetails(gid string) (*CustomerDetailRow, error) {
	url := fmt.Sprintf(
		"https://docs.google.com/spreadsheets/d/%s/export?format=csv&gid=%s",
		SPREADSHEET_ID,
		gid,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error descargando hoja individual: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error HTTP: status code %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)

	// Leer todas las filas del CSV
	allRows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error leyendo CSV: %w", err)
	}

	// Crear estructura de datos
	details := &CustomerDetailRow{}

	// Extraer datos según la estructura:
	// Fila 1 (índice 0): B1 = Nombre, E1 = Email
	// Fila 2 (índice 1): B2 = Celular, E2 = Fecha nacimiento

	// Fila 1
	if len(allRows) > 0 {
		row1 := allRows[0]
		// Columna B (índice 1)
		if len(row1) > 1 {
			details.Name = strings.TrimSpace(row1[1])
		}
		// Columna E (índice 4)
		if len(row1) > 4 {
			details.Email = strings.TrimSpace(row1[4])
		}
	}

	// Fila 2
	if len(allRows) > 1 {
		row2 := allRows[1]
		// Columna B (índice 1)
		if len(row2) > 1 {
			details.Phone = strings.TrimSpace(row2[1])
		}
		// Columna E (índice 4)
		if len(row2) > 4 {
			details.Birthday = strings.TrimSpace(row2[4])
		}
	}

	// Leer transacciones desde la fila 9 (índice 8)
	// Fila 9 es el header: Fecha, Tipo, Talla, Concepto, Tipo de pago, Valor
	// Fila 10+ son las transacciones
	details.Transactions = []TransactionRow{}
	n := allRows[0][1]
	if n == "Julián Andrés Trujillo" {
		fmt.Println("✅ Cliente encontrado: ", details.Name)
	}

	if len(allRows) > 9 { // Si hay al menos 10 filas (índice 9 = fila 10)
		for i := 9; i < len(allRows); i++ {
			row := allRows[i]

			// Validar que la fila tenga datos
			if len(row) == 0 {
				continue
			}

			// Extraer campos de la transacción
			transaction := TransactionRow{}

			// Columna A: Fecha
			if len(row) > 0 {
				transaction.Date = strings.TrimSpace(row[0])
			}

			// Columna B: Tipo (ABONO o DEUDA)
			if len(row) > 1 {
				transaction.Type = strings.TrimSpace(strings.ToUpper(row[1]))
			}

			// Columna C: Talla
			if len(row) > 2 {
				transaction.Size = strings.TrimSpace(row[2])
			}

			// Columna D: Concepto
			if len(row) > 3 {
				transaction.Concept = strings.TrimSpace(row[3])
			}

			// Columna E: Tipo de pago
			if len(row) > 4 {
				transaction.PaymentMethod = strings.TrimSpace(row[4])
			}

			// Columna F: Valor
			if len(row) > 5 {
				transaction.Amount = strings.TrimSpace(row[5])
			}

			// Validar que tenga al menos fecha y tipo
			if transaction.Date == "" || transaction.Type == "" {
				continue
			}

			// Saltar si es el header de nuevo
			if strings.EqualFold(transaction.Date, "Fecha") ||
				strings.EqualFold(transaction.Type, "Tipo") {
				continue
			}

			details.Transactions = append(details.Transactions, transaction)
		}
	}

	return details, nil
}

func connectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "fashionblue"),
		getEnv("DB_PASSWORD", "fashionblue123"),
		getEnv("DB_NAME", "fashionblue_db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	fmt.Println("✅ Conectado a la base de datos")
	return db, nil
}

func customerExists(db *gorm.DB, name, phone string) (bool, uint, error) {
	var customer models.CustomerModel

	query := db.Model(&models.CustomerModel{})

	if name != "" && phone != "" {
		query = query.Where("LOWER(name) = LOWER(?) OR phone = ?", name, phone)
	} else if name != "" {
		query = query.Where("LOWER(name) = LOWER(?)", name)
	} else if phone != "" {
		query = query.Where("phone = ?", phone)
	} else {
		return false, 0, nil
	}

	err := query.First(&customer).Error
	if err == gorm.ErrRecordNotFound {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}

	return true, customer.ID, nil
}

func createCustomer(db *gorm.DB, details *CustomerDetailRow) error {
	ctx := context.Background()

	// Parsear fecha de cumpleaños si existe
	var birthday *time.Time
	if details.Birthday != "" {
		parsed, err := parseDate(details.Birthday)
		if err == nil {
			birthday = &parsed
		} else {
			log.Printf("   ⚠️  No se pudo parsear la fecha de cumpleaños: %s", details.Birthday)
		}
	}

	// Crear entidad
	customer := &entities.Customer{
		Name:             details.Name,
		Phone:            details.Phone,
		Address:          "", // No está en la estructura actual
		RiskLevel:        entities.RiskLevelLow,
		Birthday:         birthday,
		IsActive:         true,
		PaymentFrequency: entities.PaymentFrequencyNone,
	}

	// Convertir a modelo
	model := &models.CustomerModel{}
	model.FromEntity(customer)

	// Agregar email (si el modelo lo soporta, si no, se ignora)
	// Nota: Actualmente Customer no tiene email, pero lo guardamos en el log
	if details.Email != "" {
		fmt.Printf("   📧 Email: %s (guardado en notas temporalmente)\n", details.Email)
		// Puedes agregar el email en las notas temporalmente
		customer.Notes = fmt.Sprintf("Email: %s", details.Email)
		model.FromEntity(customer)
	}

	// Crear en BD
	if err := db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	customer.ID = model.ID
	fmt.Printf("   👤 Cliente ID: %d\n", customer.ID)
	if details.Phone != "" {
		fmt.Printf("   📱 Teléfono: %s\n", details.Phone)
	}
	if birthday != nil {
		fmt.Printf("   🎂 Cumpleaños: %s\n", birthday.Format("2006-01-02"))
	}

	// Crear transacciones del cliente
	if len(details.Transactions) > 0 {
		fmt.Printf("   💰 Importando %d transacciones...\n", len(details.Transactions))

		transactionsCreated := 0
		transactionErrors := 0

		for _, txRow := range details.Transactions {
			if err := createTransaction(db, ctx, customer.ID, txRow); err != nil {
				log.Printf("      ⚠️  Error creando transacción: %v", err)
				transactionErrors++
			} else {
				transactionsCreated++
			}
		}

		if transactionsCreated > 0 {
			fmt.Printf("   ✅ %d transacciones creadas", transactionsCreated)
			if transactionErrors > 0 {
				fmt.Printf(" (%d errores)", transactionErrors)
			}
			fmt.Println()
		}
	}

	return nil
}

// createTransaction crea una transacción individual en la base de datos
func createTransaction(db *gorm.DB, ctx context.Context, customerID uint, txRow TransactionRow) error {
	// Parsear la fecha
	txDate, err := parseDate(txRow.Date)
	if err != nil {
		return fmt.Errorf("fecha inválida '%s': %w", txRow.Date, err)
	}

	// Parsear el monto
	amount, err := parseAmount(txRow.Amount)
	if err != nil {
		return fmt.Errorf("monto inválido '%s': %w", txRow.Amount, err)
	}

	// Determinar el tipo de transacción
	var txType entities.TransactionType
	switch strings.ToUpper(txRow.Type) {
	case "DEUDA":
		txType = entities.TransactionTypeDebt
	case "ABONO":
		txType = entities.TransactionTypePayment
	default:
		return fmt.Errorf("tipo de transacción inválido: %s", txRow.Type)
	}

	// Construir descripción
	description := txRow.Concept
	if txRow.Size != "" {
		description = fmt.Sprintf("%s (Talla: %s)", description, txRow.Size)
	}
	if description == "" {
		description = "-"
	}

	// Buscar método de pago si existe
	var paymentMethodID *uint
	if txRow.PaymentMethod != "" && txType == entities.TransactionTypePayment {
		// Buscar método de pago en la BD
		var paymentMethod models.PaymentMethodModel
		if err := db.Where("LOWER(name) = LOWER(?)", txRow.PaymentMethod).First(&paymentMethod).Error; err == nil {
			paymentMethodID = &paymentMethod.ID
		}
	}

	// Crear transacción
	transaction := &entities.CustomerTransaction{
		CustomerID:      customerID,
		Type:            txType,
		Amount:          amount,
		Description:     description,
		PaymentMethodID: paymentMethodID,
		Date:            txDate,
	}

	txModel := &models.CustomerTransactionModel{}
	txModel.FromEntity(transaction)

	if err := db.WithContext(ctx).Create(txModel).Error; err != nil {
		return err
	}

	return nil
}

// parseAmount parsea un monto de dinero (elimina $, comas, espacios)
func parseAmount(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, " ", "")

	if s == "" {
		return 0, fmt.Errorf("monto vacío")
	}

	// Convertir a float
	amount, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

// parseDate intenta parsear una fecha en diferentes formatos
func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
		"02-01-2006",
		"01-02-2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("formato de fecha no reconocido: %s", dateStr)
}

func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, " ", "")

	if s == "" {
		return 0, nil
	}

	return strconv.ParseFloat(s, 64)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
