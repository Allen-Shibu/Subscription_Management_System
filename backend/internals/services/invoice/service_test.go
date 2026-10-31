package invoice

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"recur_backend/internals/models"
	"recur_backend/pkg/email"
)

type mockRepo struct {
	getInvoiceDetailsFunc    func(ctx context.Context, id int64) (models.InvoiceDetails, error)
	markInvoiceEmailSentFunc func(ctx context.Context, id int64) (models.Invoice, error)
}

func (m *mockRepo) GetInvoiceDetails(ctx context.Context, id int64) (models.InvoiceDetails, error) {
	if m.getInvoiceDetailsFunc != nil {
		return m.getInvoiceDetailsFunc(ctx, id)
	}
	return models.InvoiceDetails{}, errors.New("not implemented")
}

func (m *mockRepo) MarkInvoiceEmailSent(ctx context.Context, id int64) (models.Invoice, error) {
	if m.markInvoiceEmailSentFunc != nil {
		return m.markInvoiceEmailSentFunc(ctx, id)
	}
	return models.Invoice{}, errors.New("not implemented")
}

type mockSender struct {
	sent     []email.Message
	sendFunc func(ctx context.Context, msg email.Message) error
}

func (m *mockSender) Send(ctx context.Context, msg email.Message) error {
	m.sent = append(m.sent, msg)
	if m.sendFunc != nil {
		return m.sendFunc(ctx, msg)
	}
	return nil
}

func sampleDetails() models.InvoiceDetails {
	issued := time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC)
	due := issued.AddDate(0, 0, 14)
	return models.InvoiceDetails{
		Invoice: models.Invoice{
			ID:            101,
			InvoiceNumber: "INV-000101",
			Status:        models.InvoiceStatusOpen,
			PeriodStart:   issued,
			PeriodEnd:     issued.AddDate(0, 1, 0),
			SubtotalCents: 5000,
			DiscountCents: 500,
			TaxCents:      810,
			TotalCents:    5310,
			TaxRateBps:    1800,
			IssuedAt:      &issued,
			DueAt:         &due,
		},
		CustomerName:    "Alice Smith",
		CustomerEmail:   "alice@example.com",
		ProductName:     "Cloud Pro",
		VariantName:     "Annual",
		Quantity:        1,
		UnitPriceCents:  5000,
		Currency:        "USD",
		BillingInterval: "year",
	}
}

func TestRender_Success(t *testing.T) {
	details := sampleDetails()
	repo := &mockRepo{
		getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
			if id != 101 {
				t.Fatalf("expected id 101, got %d", id)
			}
			return details, nil
		},
	}
	svc := New(repo, &mockSender{})

	doc, gotDetails, err := svc.Render(context.Background(), 101)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.HasPrefix(doc, []byte("%PDF-")) {
		t.Fatalf("expected PDF output, got %q", doc[:10])
	}
	if gotDetails.InvoiceNumber != "INV-000101" {
		t.Fatalf("expected INV-000101, got %s", gotDetails.InvoiceNumber)
	}
}

func TestRender_RepoError(t *testing.T) {
	expectedErr := errors.New("database connection failed")
	repo := &mockRepo{
		getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
			return models.InvoiceDetails{}, expectedErr
		},
	}
	svc := New(repo, &mockSender{})

	_, _, err := svc.Render(context.Background(), 1)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestSend_Success_OpenInvoice(t *testing.T) {
	details := sampleDetails()
	repoCalled := false
	repo := &mockRepo{
		getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
			return details, nil
		},
		markInvoiceEmailSentFunc: func(ctx context.Context, id int64) (models.Invoice, error) {
			repoCalled = true
			now := time.Now().UTC()
			inv := details.Invoice
			inv.EmailSentAt = &now
			return inv, nil
		},
	}
	sender := &mockSender{}
	svc := New(repo, sender)

	inv, err := svc.Send(context.Background(), 101)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !repoCalled {
		t.Fatal("MarkInvoiceEmailSent was not called")
	}
	if inv.EmailSentAt == nil {
		t.Fatal("EmailSentAt should be populated")
	}
	if len(sender.sent) != 1 {
		t.Fatalf("expected 1 sent message, got %d", len(sender.sent))
	}

	msg := sender.sent[0]
	if msg.To != "alice@example.com" {
		t.Errorf("expected to alice@example.com, got %s", msg.To)
	}
	if !strings.Contains(msg.Subject, "INV-000101") {
		t.Errorf("subject should contain invoice number, got %s", msg.Subject)
	}
	if msg.AttachmentName != "INV-000101.pdf" {
		t.Errorf("expected attachment name INV-000101.pdf, got %s", msg.AttachmentName)
	}
	if !bytes.HasPrefix(msg.Attachment, []byte("%PDF-")) {
		t.Errorf("attachment is not a valid PDF")
	}
	if !strings.Contains(msg.Text, "Total due: USD 53.10") {
		t.Errorf("text body missing total due: %s", msg.Text)
	}
	if !strings.Contains(msg.HTML, "Total due") {
		t.Errorf("html body missing total due: %s", msg.HTML)
	}
}

func TestSend_Success_PaidInvoice(t *testing.T) {
	details := sampleDetails()
	details.Status = models.InvoiceStatusPaid
	now := time.Now().UTC()
	details.PaidAt = &now

	repo := &mockRepo{
		getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
			return details, nil
		},
		markInvoiceEmailSentFunc: func(ctx context.Context, id int64) (models.Invoice, error) {
			inv := details.Invoice
			inv.EmailSentAt = &now
			return inv, nil
		},
	}
	sender := &mockSender{}
	svc := New(repo, sender)

	_, err := svc.Send(context.Background(), 101)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sender.sent) != 1 {
		t.Fatalf("expected 1 sent message, got %d", len(sender.sent))
	}
	msg := sender.sent[0]
	if !strings.Contains(msg.Text, "Total paid: USD 53.10") {
		t.Errorf("text body should mention total paid: %s", msg.Text)
	}
	if !strings.Contains(msg.HTML, "Total paid") {
		t.Errorf("html body should mention total paid: %s", msg.HTML)
	}
}

func TestSend_DisallowedStatuses(t *testing.T) {
	statuses := []string{
		models.InvoiceStatusDraft,
		models.InvoiceStatusVoid,
		models.InvoiceStatusUncollectible,
	}

	for _, st := range statuses {
		t.Run("status_"+st, func(t *testing.T) {
			details := sampleDetails()
			details.Status = st

			repo := &mockRepo{
				getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
					return details, nil
				},
			}
			sender := &mockSender{}
			svc := New(repo, sender)

			_, err := svc.Send(context.Background(), 101)
			if !errors.Is(err, ErrInvoiceNotSendable) {
				t.Fatalf("expected ErrInvoiceNotSendable, got %v", err)
			}
			if len(sender.sent) > 0 {
				t.Fatal("no email should be sent for disallowed status")
			}
		})
	}
}

func TestSend_AlreadySent(t *testing.T) {
	details := sampleDetails()
	now := time.Now().UTC()
	details.EmailSentAt = &now

	repo := &mockRepo{
		getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
			return details, nil
		},
	}
	sender := &mockSender{}
	svc := New(repo, sender)

	_, err := svc.Send(context.Background(), 101)
	if !errors.Is(err, ErrInvoiceAlreadySent) {
		t.Fatalf("expected ErrInvoiceAlreadySent, got %v", err)
	}
	if len(sender.sent) > 0 {
		t.Fatal("no email should be sent")
	}
}

func TestSend_MissingEmail(t *testing.T) {
	details := sampleDetails()
	details.CustomerEmail = "   "

	repo := &mockRepo{
		getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
			return details, nil
		},
	}
	sender := &mockSender{}
	svc := New(repo, sender)

	_, err := svc.Send(context.Background(), 101)
	if !errors.Is(err, ErrCustomerEmailMissing) {
		t.Fatalf("expected ErrCustomerEmailMissing, got %v", err)
	}
	if len(sender.sent) > 0 {
		t.Fatal("no email should be sent")
	}
}

func TestSend_EmailDeliveryFailure(t *testing.T) {
	details := sampleDetails()
	repoCalled := false
	repo := &mockRepo{
		getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
			return details, nil
		},
		markInvoiceEmailSentFunc: func(ctx context.Context, id int64) (models.Invoice, error) {
			repoCalled = true
			return details.Invoice, nil
		},
	}
	sender := &mockSender{
		sendFunc: func(ctx context.Context, msg email.Message) error {
			return errors.New("smtp timeout")
		},
	}
	svc := New(repo, sender)

	_, err := svc.Send(context.Background(), 101)
	if !errors.Is(err, ErrEmailDelivery) {
		t.Fatalf("expected ErrEmailDelivery, got %v", err)
	}
	if repoCalled {
		t.Fatal("MarkInvoiceEmailSent must NOT be called when delivery fails")
	}
}
