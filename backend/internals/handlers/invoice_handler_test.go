package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"recur_backend/internals/models"
	invoicesvc "recur_backend/internals/services/invoice"
	"recur_backend/pkg/email"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type mockInvoiceRepo struct {
	getInvoiceDetailsFunc    func(ctx context.Context, id int64) (models.InvoiceDetails, error)
	markInvoiceEmailSentFunc func(ctx context.Context, id int64) (models.Invoice, error)
}

func (m *mockInvoiceRepo) GetInvoiceDetails(ctx context.Context, id int64) (models.InvoiceDetails, error) {
	if m.getInvoiceDetailsFunc != nil {
		return m.getInvoiceDetailsFunc(ctx, id)
	}
	return models.InvoiceDetails{}, errors.New("not implemented")
}

func (m *mockInvoiceRepo) MarkInvoiceEmailSent(ctx context.Context, id int64) (models.Invoice, error) {
	if m.markInvoiceEmailSentFunc != nil {
		return m.markInvoiceEmailSentFunc(ctx, id)
	}
	return models.Invoice{}, errors.New("not implemented")
}

func setupInvoiceTestRouter(svc *invoicesvc.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := &SubscriptionHandler{invoices: svc}
	r.GET("/invoices/:invoiceID/pdf", h.DownloadInvoicePDF)
	r.POST("/invoices/:invoiceID/send", h.SendInvoice)
	return r
}

func testInvoiceDetails(id int64) models.InvoiceDetails {
	now := time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)
	due := now.AddDate(0, 0, 14)
	return models.InvoiceDetails{
		Invoice: models.Invoice{
			ID:            id,
			InvoiceNumber: "INV-000042",
			Status:        models.InvoiceStatusOpen,
			PeriodStart:   now,
			PeriodEnd:     now.AddDate(0, 1, 0),
			SubtotalCents: 2000,
			TotalCents:    2000,
			IssuedAt:      &now,
			DueAt:         &due,
		},
		CustomerName:    "Test User",
		CustomerEmail:   "user@example.com",
		ProductName:     "Basic Plan",
		Quantity:        1,
		UnitPriceCents:  2000,
		Currency:        "USD",
		BillingInterval: "month",
	}
}

func TestDownloadInvoicePDF_Success(t *testing.T) {
	repo := &mockInvoiceRepo{
		getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
			return testInvoiceDetails(id), nil
		},
	}
	svc := invoicesvc.New(repo, &email.NoopSender{})
	router := setupInvoiceTestRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/invoices/42/pdf", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/pdf" {
		t.Fatalf("expected application/pdf, got %s", ct)
	}
	if cd := w.Header().Get("Content-Disposition"); cd != `attachment; filename="INV-000042.pdf"` {
		t.Fatalf("expected attachment header, got %s", cd)
	}
	if len(w.Body.Bytes()) < 100 {
		t.Fatalf("suspiciously small PDF body: %d bytes", len(w.Body.Bytes()))
	}
}

func TestDownloadInvoicePDF_NotFound(t *testing.T) {
	repo := &mockInvoiceRepo{
		getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
			return models.InvoiceDetails{}, pgx.ErrNoRows
		},
	}
	svc := invoicesvc.New(repo, &email.NoopSender{})
	router := setupInvoiceTestRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/invoices/999/pdf", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d", w.Code)
	}
}

func TestDownloadInvoicePDF_InvalidID(t *testing.T) {
	svc := invoicesvc.New(&mockInvoiceRepo{}, &email.NoopSender{})
	router := setupInvoiceTestRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/invoices/invalid/pdf", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
}

func TestSendInvoice_Success(t *testing.T) {
	details := testInvoiceDetails(42)
	repo := &mockInvoiceRepo{
		getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
			return details, nil
		},
		markInvoiceEmailSentFunc: func(ctx context.Context, id int64) (models.Invoice, error) {
			now := time.Now().UTC()
			inv := details.Invoice
			inv.EmailSentAt = &now
			return inv, nil
		},
	}
	svc := invoicesvc.New(repo, &email.NoopSender{})
	router := setupInvoiceTestRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/invoices/42/send", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSendInvoice_ErrorMappings(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *invoicesvc.Service
		wantStatus int
	}{
		{
			name: "already sent",
			setup: func() *invoicesvc.Service {
				d := testInvoiceDetails(42)
				now := time.Now().UTC()
				d.EmailSentAt = &now
				repo := &mockInvoiceRepo{
					getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
						return d, nil
					},
				}
				return invoicesvc.New(repo, &email.NoopSender{})
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "void invoice",
			setup: func() *invoicesvc.Service {
				d := testInvoiceDetails(42)
				d.Status = models.InvoiceStatusVoid
				repo := &mockInvoiceRepo{
					getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
						return d, nil
					},
				}
				return invoicesvc.New(repo, &email.NoopSender{})
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "missing customer email",
			setup: func() *invoicesvc.Service {
				d := testInvoiceDetails(42)
				d.CustomerEmail = ""
				repo := &mockInvoiceRepo{
					getInvoiceDetailsFunc: func(ctx context.Context, id int64) (models.InvoiceDetails, error) {
						return d, nil
					},
				}
				return invoicesvc.New(repo, &email.NoopSender{})
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.setup()
			router := setupInvoiceTestRouter(svc)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/invoices/42/send", nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}
