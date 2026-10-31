// Package invoice renders and emails subscription invoices.
package invoice

import (
	"context"
	"errors"
	"fmt"
	"html"
	"strings"

	"recur_backend/internals/models"
	"recur_backend/pkg/email"
	"recur_backend/pkg/pdf"
)

var (
	ErrInvoiceAlreadySent    = errors.New("invoice has already been emailed")
	ErrCustomerEmailMissing  = errors.New("customer does not have an email address")
	ErrInvoiceNotSendable    = errors.New("invoice cannot be emailed in its current status")
	ErrEmailDelivery         = errors.New("failed to deliver invoice email")
)

// Repository defines the subscription and invoice operations needed by Service.
type Repository interface {
	GetInvoiceDetails(ctx context.Context, id int64) (models.InvoiceDetails, error)
	MarkInvoiceEmailSent(ctx context.Context, id int64) (models.Invoice, error)
}

// Service turns invoices into PDFs and sends them to customers.
type Service struct {
	repo   Repository
	sender email.Sender
}

func New(repo Repository, sender email.Sender) *Service {
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
	if details.Status == models.InvoiceStatusDraft ||
		details.Status == models.InvoiceStatusVoid ||
		details.Status == models.InvoiceStatusUncollectible {
		return models.Invoice{}, ErrInvoiceNotSendable
	}
	if details.EmailSentAt != nil {
		return models.Invoice{}, ErrInvoiceAlreadySent
	}
	customerEmail := strings.TrimSpace(details.CustomerEmail)
	if customerEmail == "" {
		return models.Invoice{}, ErrCustomerEmailMissing
	}
	details.CustomerEmail = customerEmail

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
	fmt.Fprintf(&b, "Status: %s\n", strings.ToUpper(d.Status))
	if d.ProductName != "" {
		item := d.ProductName
		if d.VariantName != "" && d.VariantName != d.ProductName {
			item += " - " + d.VariantName
		}
		fmt.Fprintf(&b, "Plan: %s\n", item)
	}
	fmt.Fprintf(&b, "Period: %s to %s\n", d.PeriodStart.Format("02 Jan 2006"), d.PeriodEnd.Format("02 Jan 2006"))
	fmt.Fprintf(&b, "Subtotal: %s\n", money(d.SubtotalCents, d.Currency))
	if d.DiscountCents > 0 {
		fmt.Fprintf(&b, "Discount: -%s\n", money(d.DiscountCents, d.Currency))
	}
	if d.TaxCents > 0 {
		taxLabel := "Tax"
		if d.TaxRateBps > 0 {
			if d.TaxRateBps%100 == 0 {
				taxLabel = fmt.Sprintf("Tax (%d%%)", d.TaxRateBps/100)
			} else {
				taxLabel = fmt.Sprintf("Tax (%.2f%%)", float64(d.TaxRateBps)/100.0)
			}
		}
		fmt.Fprintf(&b, "%s: %s\n", taxLabel, money(d.TaxCents, d.Currency))
	}
	if d.Status == models.InvoiceStatusPaid {
		fmt.Fprintf(&b, "Total paid: %s\n", money(d.TotalCents, d.Currency))
		if d.PaidAt != nil {
			fmt.Fprintf(&b, "Paid date: %s\n", d.PaidAt.Format("02 Jan 2006"))
		}
	} else {
		fmt.Fprintf(&b, "Total due: %s\n", money(d.TotalCents, d.Currency))
		if d.DueAt != nil {
			fmt.Fprintf(&b, "Due date: %s\n", d.DueAt.Format("02 Jan 2006"))
		}
	}
	b.WriteString("\nThank you for your business.\n\nBrightloop Software\n")
	return b.String()
}

func htmlBody(d models.InvoiceDetails) string {
	name := html.EscapeString(d.CustomerName)
	number := html.EscapeString(d.InvoiceNumber)
	status := strings.ToUpper(d.Status)

	var rows strings.Builder
	if d.ProductName != "" {
		item := html.EscapeString(d.ProductName)
		if d.VariantName != "" && d.VariantName != d.ProductName {
			item += " - " + html.EscapeString(d.VariantName)
		}
		fmt.Fprintf(&rows, `<tr><td style="padding:6px 0;color:#71717a">Plan</td><td style="padding:6px 0;text-align:right;font-weight:500">%s</td></tr>`, item)
	}
	fmt.Fprintf(&rows, `<tr><td style="padding:6px 0;color:#71717a">Period</td><td style="padding:6px 0;text-align:right">%s to %s</td></tr>`,
		d.PeriodStart.Format("02 Jan 2006"), d.PeriodEnd.Format("02 Jan 2006"))
	fmt.Fprintf(&rows, `<tr><td style="padding:6px 0;color:#71717a">Subtotal</td><td style="padding:6px 0;text-align:right">%s</td></tr>`,
		money(d.SubtotalCents, d.Currency))
	if d.DiscountCents > 0 {
		fmt.Fprintf(&rows, `<tr><td style="padding:6px 0;color:#71717a">Discount</td><td style="padding:6px 0;text-align:right;color:#16a34a">-%s</td></tr>`,
			money(d.DiscountCents, d.Currency))
	}
	if d.TaxCents > 0 {
		taxLabel := "Tax"
		if d.TaxRateBps > 0 {
			if d.TaxRateBps%100 == 0 {
				taxLabel = fmt.Sprintf("Tax (%d%%)", d.TaxRateBps/100)
			} else {
				taxLabel = fmt.Sprintf("Tax (%.2f%%)", float64(d.TaxRateBps)/100.0)
			}
		}
		fmt.Fprintf(&rows, `<tr><td style="padding:6px 0;color:#71717a">%s</td><td style="padding:6px 0;text-align:right">%s</td></tr>`,
			html.EscapeString(taxLabel), money(d.TaxCents, d.Currency))
	}

	if d.Status == models.InvoiceStatusPaid {
		if d.PaidAt != nil {
			fmt.Fprintf(&rows, `<tr><td style="padding:6px 0;color:#71717a">Paid on</td><td style="padding:6px 0;text-align:right">%s</td></tr>`,
				d.PaidAt.Format("02 Jan 2006"))
		}
		fmt.Fprintf(&rows, `<tr style="border-top:1px solid #e4e4e7"><td style="padding:10px 0;font-weight:bold">Total paid</td><td style="padding:10px 0;text-align:right;font-weight:bold">%s</td></tr>`,
			money(d.TotalCents, d.Currency))
	} else {
		due := "on receipt"
		if d.DueAt != nil {
			due = d.DueAt.Format("02 Jan 2006")
		}
		fmt.Fprintf(&rows, `<tr><td style="padding:6px 0;color:#71717a">Due date</td><td style="padding:6px 0;text-align:right">%s</td></tr>`, due)
		fmt.Fprintf(&rows, `<tr style="border-top:1px solid #e4e4e7"><td style="padding:10px 0;font-weight:bold">Total due</td><td style="padding:10px 0;text-align:right;font-weight:bold">%s</td></tr>`,
			money(d.TotalCents, d.Currency))
	}

	return fmt.Sprintf(`<div style="font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#18181b;max-width:560px;margin:0 auto;padding:24px 16px;line-height:1.5">
<h2 style="margin:0 0 16px;font-size:20px;font-weight:600">Invoice %s</h2>
<p style="margin:0 0 16px">Hello %s,</p>
<p style="margin:0 0 20px">Your invoice <strong>%s</strong> (Status: <strong>%s</strong>) is attached to this email as a PDF document.</p>
<table style="width:100%%;border-collapse:collapse;margin:0 0 24px">
%s
</table>
<p style="margin:0 0 4px">Thank you for your business.</p>
<p style="margin:0;color:#71717a;font-size:14px">Brightloop Software</p>
</div>`,
		number, name, number, status, rows.String())
}

func money(cents int64, currency string) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s%s %d.%02d", sign, currency, cents/100, cents%100)
}
