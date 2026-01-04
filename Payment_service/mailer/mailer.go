package mailer

import (
	"bytes"
	"fmt"
	"log"

	"github.com/go-mail/mail"
	"github.com/jung-kurt/gofpdf"
)

type Mailer struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewMailer(host string, port int, username, password, from string) *Mailer {
	return &Mailer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

func (m *Mailer) Send(to string, subject string, plainBody string) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plainBody)

	d := mail.NewDialer(m.Host, m.Port, m.Username, m.Password)
	return d.DialAndSend(msg)
}

func (m *Mailer) SendWithPDF(to string, subject string, htmlBody string, pdfBytes []byte) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)

	msg.AttachReader("receipt.pdf", bytes.NewReader(pdfBytes))

	d := mail.NewDialer(m.Host, m.Port, m.Username, m.Password)
	return d.DialAndSend(msg)
}

func (m *Mailer) GeneratePDFReceipt(orderID, userID string, items []string, total float64) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 14)

	pdf.Cell(40, 10, fmt.Sprintf("Order Receipt #%s", orderID))
	pdf.Ln(12)
	pdf.Cell(40, 10, fmt.Sprintf("User ID: %s", userID))
	pdf.Ln(12)

	pdf.Cell(40, 10, "Items:")
	for _, item := range items {
		pdf.Ln(8)
		pdf.Cell(40, 10, fmt.Sprintf("- %s", item))
	}

	pdf.Ln(12)
	pdf.Cell(40, 10, fmt.Sprintf("Total: %.2f USD", total))

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		log.Printf("PDF generation failed: %v", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
