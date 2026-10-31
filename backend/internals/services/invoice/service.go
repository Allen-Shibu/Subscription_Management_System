// Package invoice renders and emails subscription invoices.
package invoice

import (
	"context"
	"errors"
	"fmt"
	"html"
	"strings"

	"recur_backend/internals/models"
	"recur_backend/internals/repository"
	"recur_backend/pkg/email"
	"recur_backend/pkg/pdf"
)

var (
	ErrInvoiceAlreadySent    = errors.New("invoice has already been emailed")
	ErrCustomerEmailMissing  = errors.New("customer does not have an email address")
	ErrInvoiceNotSendable    = errors.New("invoice cannot be emailed in its current status")
	ErrEmailDelivery         = errors.New("failed to deliver invoice email")
)

// Service turns invoices into PDFs and sends them to customers.
type Service struct {
	repo   *repository.SubscriptionRepository
	sender email.Sender
}

func New(repo *repository.SubscriptionRepository, sender email.Sender) *Service {
	return &Service{repo: repo, sender: sender}
}

// Render returns the invoice PDF and the details used to build it.
func (s *Service) Render(ctx context.Context, invoiceID int64) ([]byte, models.InvoiceDetails, error) {
	details, err := s.repo.GetInvoiceDetails(ctx, invoiceID)
	if err != nil {
		return nil, models.InvoiceDetails{}, err
	}
	document, err := pdf.RenderInvoice(details)
	if err != nil {
		return nil, models.InvoiceDetails{}, err
	}
	return document, details, nil
}

// Send renders the invoice and emails it as a PDF attachment, then records
// the delivery time. Invoices can only be sent once.
func (s *Service) Send(ctx context.Context, invoiceID int64) (models.Invoice, error) {
	details, err := s.repo.GetInvoiceDetails(ctx, invoiceID)
	if err != nil {
		return models.Invoice{}, err
	}
	if details.Status == models.InvoiceStatusVoid ||
		details.Status == models.InvoiceStatusUncollectible {
		return models.Invoice{}, ErrInvoiceNotSendable
	}
	if details.EmailSentAt != nil {
		return models.Invoice{}, ErrInvoiceAlreadySent
	}
	if details.CustomerEmail == "" {
		return models.Invoice{}, ErrCustomerEmailMissing
	}

	document, err := pdf.RenderInvoice(details)
	if err != nil {
		return models.Invoice{}, err
	}

	msg := email.Message{
		To:             details.CustomerEmail,
		Subject:        "Invoice " + details.InvoiceNumber + " from Brightloop Software",
		Text:           textBody(details),
		HTML:           htmlBody(details),
		AttachmentName: details.InvoiceNumber + ".pdf",
		Attachment:     document,
	}
	if err := s.sender.Send(ctx, msg); err != nil {
		return models.Invoice{}, fmt.Errorf("%w: %v", ErrEmailDelivery, err)
	}

	return s.repo.MarkInvoiceEmailSent(ctx, invoiceID)
}

func textBody(d models.InvoiceDetails) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Hello %s,\n\n", d.CustomerName)
	fmt.Fprintf(&b, "Your invoice %s is attached.\n\n", d.InvoiceNumber)
	fmt.Fprintf(&b, "Period: %s to %s\n", d.PeriodStart.Format("02 Jan 2006"), d.PeriodEnd.Format("02 Jan 2006"))
	fmt.Fprintf(&b, "Total due: %s\n", money(d.TotalCents, d.Currency))
	if d.DueAt != nil {
		fmt.Fprintf(&b, "Due date: %s\n", d.DueAt.Format("02 Jan 2006"))
	}
	b.WriteString("\nThank you for your business.\n\nBrightloop Software\n")
	return b.String()
}

func htmlBody(d models.InvoiceDetails) string {
	name := html.EscapeString(d.CustomerName)
	number := html.EscapeString(d.InvoiceNumber)
	due := "on receipt"
	if d.DueAt != nil {
		due = d.DueAt.Format("02 Jan 2006")
	}
	return fmt.Sprintf(`<p>Hello %s,</p>
<p>Your invoice <strong>%s</strong> is attached to this email.</p>
<table cellpadding="4" style="border-collapse:collapse">
  <tr><td>Period</td><td>%s to %s</td></tr>
  <tr><td>Due date</td><td>%s</td></tr>
  <tr><td><strong>Total due</strong></td><td><strong>%s</strong></td></tr>
</table>
<p>Thank you for your business.<br>Brightloop Software</p>`,
		name, number,
		d.PeriodStart.Format("02 Jan 2006"), d.PeriodEnd.Format("02 Jan 2006"),
		due, money(d.TotalCents, d.Currency))
}

func money(cents int64, currency string) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s%s %d.%02d", sign, currency, cents/100, cents%100)
}
