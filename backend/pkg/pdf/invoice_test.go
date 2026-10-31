package pdf

import (
	"bytes"
	"testing"
	"time"

	"recur_backend/internals/models"
)

func TestRenderInvoice(t *testing.T) {
	issued := time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC)
	due := issued.AddDate(0, 0, 14)
	details := models.InvoiceDetails{
		Invoice: models.Invoice{
			ID:            1,
			InvoiceNumber: "INV-000001",
			Status:        models.InvoiceStatusOpen,
			PeriodStart:   issued,
			PeriodEnd:     issued.AddDate(0, 1, 0),
			SubtotalCents: 199900,
			DiscountCents: 10000,
			TotalCents:    189900,
			IssuedAt:      &issued,
			DueAt:         &due,
		},
		CustomerName:    "Ada Lovelace",
		CustomerEmail:   "ada@example.com",
		ProductName:     "Cloud Suite",
		VariantName:     "Growth",
		Quantity:        2,
		UnitPriceCents:  99950,
		Currency:        "USD",
		BillingInterval: "month",
	}

	document, err := RenderInvoice(details)
	if err != nil {
		t.Fatalf("RenderInvoice() error = %v", err)
	}
	if !bytes.HasPrefix(document, []byte("%PDF-")) {
		t.Fatalf("document is not a PDF")
	}
	if len(document) < 1000 {
		t.Fatalf("document is suspiciously small: %d bytes", len(document))
	}
}

func TestFormatMoney(t *testing.T) {
	tests := []struct {
		cents    int64
		currency string
		want     string
	}{
		{0, "USD", "USD 0.00"},
		{1999, "USD", "USD 19.99"},
		{123456789, "EUR", "EUR 1,234,567.89"},
		{-500, "USD", "-USD 5.00"},
	}
	for _, tt := range tests {
		if got := formatMoney(tt.cents, tt.currency); got != tt.want {
			t.Fatalf("formatMoney(%d, %q) = %q, want %q", tt.cents, tt.currency, got, tt.want)
		}
	}
}
