package nats

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"payment/mailer"
)

type OrderCreatedEvent struct {
	OrderID   string   `json:"orderId"`
	UserID    string   `json:"userId"`
	Items     []string `json:"items"`
	Total     float64  `json:"total"`
	CreatedAt string   `json:"createdAt"`
}

type EmailWorker struct {
	Mailer     *mailer.Mailer
	GetEmailFn func(userID string) (string, error)
}

func (e *EmailWorker) HandleOrderCreated(m *nats.Msg) {
	var evt OrderCreatedEvent
	if err := json.Unmarshal(m.Data, &evt); err != nil {
		log.Printf("[NATS] Invalid event: %v", err)
		return
	}

	log.Printf("[NATS] Received order.created: %s", evt.OrderID)

	email, err := e.GetEmailFn(evt.UserID)
	if err != nil {
		log.Printf("[EMAIL]Failed to get email for user %s: %v", evt.UserID, err)
		return
	}

	html := generateHTML(evt)
	pdf, _ := e.Mailer.GeneratePDFReceipt(evt.OrderID, evt.UserID, evt.Items, evt.Total)

	err = e.Mailer.SendWithPDF(email, "Order Receipt", html, pdf)
	if err != nil {
		log.Printf("[EMAIL] Failed to send: %v", err)
	} else {
		log.Printf("[EMAIL] Receipt sent to %s", email)
	}
}

func generateHTML(evt OrderCreatedEvent) string {
	list := ""
	for _, item := range evt.Items {
		list += "<li>" + item + "</li>"
	}
	return `
		<h2>Order Receipt</h2>
		<p><strong>Order ID:</strong> ` + evt.OrderID + `</p>
		<p><strong>Total:</strong> $` + formatPrice(evt.Total) + `</p>
		<p><strong>Created At:</strong> ` + evt.CreatedAt + `</p>
		<ul>` + list + `</ul>
	`
}

func formatPrice(f float64) string {
	return fmt.Sprintf("%.2f", f)
}
