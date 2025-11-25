package event_handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
)

// WebhookConfig configuración de webhooks
type WebhookConfig struct {
	URL     string
	Enabled bool
	Secret  string // Para firmar requests
}

// WebhookHandler maneja eventos para enviar webhooks
type WebhookHandler struct {
	eventBus  *events.EventBus
	eventChan chan events.OrderEvent
	stopChan  chan bool
	config    WebhookConfig
	client    *http.Client
}

// NewWebhookHandler crea un nuevo handler de webhooks
func NewWebhookHandler(eventBus *events.EventBus, config WebhookConfig) *WebhookHandler {
	handler := &WebhookHandler{
		eventBus:  eventBus,
		eventChan: make(chan events.OrderEvent, 100),
		stopChan:  make(chan bool),
		config:    config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	// Suscribirse a eventos importantes para webhooks
	eventBus.Subscribe(events.EventOrderApproved, handler.eventChan)
	eventBus.Subscribe(events.EventOrderDelivered, handler.eventChan)
	eventBus.Subscribe(events.EventOrderCancelled, handler.eventChan)
	eventBus.Subscribe(events.EventStockUpdated, handler.eventChan)

	return handler
}

// Start inicia el procesamiento de eventos
func (h *WebhookHandler) Start() {
	if !h.config.Enabled {
		log.Println("🔗 Webhook Event Handler disabled")
		return
	}

	log.Println("🔗 Webhook Event Handler started")

	go func() {
		for {
			select {
			case event := <-h.eventChan:
				h.handleEvent(event)
			case <-h.stopChan:
				log.Println("🔗 Webhook Event Handler stopped")
				return
			}
		}
	}()
}

// Stop detiene el procesamiento de eventos
func (h *WebhookHandler) Stop() {
	h.stopChan <- true
}

// handleEvent procesa un evento y envía webhook
func (h *WebhookHandler) handleEvent(event events.OrderEvent) {
	if !h.config.Enabled || h.config.URL == "" {
		return
	}

	// Preparar payload
	payload := map[string]interface{}{
		"event_type": event.Type,
		"order_id":   event.OrderID,
		"old_status": event.OldStatus,
		"new_status": event.NewStatus,
		"timestamp":  event.Timestamp,
		"data":       event.Data,
	}

	// Enviar webhook
	if err := h.sendWebhook(payload); err != nil {
		log.Printf("🔗 [WEBHOOK ERROR] Failed to send webhook for order #%d: %v", event.OrderID, err)
	} else {
		log.Printf("🔗 [WEBHOOK] Sent webhook for order #%d - Event: %s", event.OrderID, event.Type)
	}
}

// sendWebhook envía el webhook al endpoint configurado
func (h *WebhookHandler) sendWebhook(payload map[string]interface{}) error {
	// Serializar payload
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Crear request
	req, err := http.NewRequest("POST", h.config.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "FashionBlue-Webhook/1.0")

	// Agregar firma si hay secret configurado
	if h.config.Secret != "" {
		// TODO: Implementar firma HMAC
		// signature := generateHMAC(jsonData, h.config.Secret)
		// req.Header.Set("X-Webhook-Signature", signature)
	}

	// Enviar request
	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verificar respuesta
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("🔗 [WEBHOOK] Received non-2xx status code: %d", resp.StatusCode)
	}

	return nil
}
