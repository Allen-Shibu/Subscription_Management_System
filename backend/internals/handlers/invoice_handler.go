package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"recur_backend/internals/middleware"
	"recur_backend/internals/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type refundInput struct {
	AmountCents int64  `json:"amountCents"`
	Reason      string `json:"reason"`
}

func (in *refundInput) validate() error {
	in.Reason = strings.TrimSpace(in.Reason)
	if in.AmountCents <= 0 {
		return errors.New("amountCents must be greater than zero")
	}
	if in.Reason == "" {
		return errors.New("reason is required")
	}
	return nil
}

// CreateInvoice issues an invoice for the subscription's current period.
func (h *SubscriptionHandler) CreateInvoice(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	subscription, err := h.repo.GetSubscription(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}
	if subscription.Status == models.SubscriptionStatusCanceled ||
		subscription.Status == models.SubscriptionStatusExpired {
		c.JSON(http.StatusConflict, gin.H{"error": "subscription is not billable"})
		return
	}
	createdBy, _ := middleware.GetUserID(c)

	invoice, err := h.repo.CreateInvoice(c.Request.Context(), id, createdBy)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusCreated, invoice)
}

func (h *SubscriptionHandler) ListInvoices(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	invoices, err := h.repo.ListInvoices(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, invoices)
}

func (h *SubscriptionHandler) GetInvoice(c *gin.Context) {
	id, ok := parseID(c, "invoiceID")
	if !ok {
		return
	}

	invoice, err := h.repo.GetInvoice(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, invoice)
}

func (h *SubscriptionHandler) PayInvoice(c *gin.Context) {
	id, ok := parseID(c, "invoiceID")
	if !ok {
		return
	}

	invoice, err := h.repo.GetInvoice(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}
	if invoice.Status != models.InvoiceStatusOpen {
		c.JSON(http.StatusConflict, gin.H{"error": "only open invoices can be paid"})
		return
	}

	paid, err := h.repo.MarkInvoicePaid(c.Request.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusConflict, gin.H{"error": "invoice is no longer open"})
		return
	}
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, paid)
}

func (h *SubscriptionHandler) CreateRefund(c *gin.Context) {
	id, ok := parseID(c, "invoiceID")
	if !ok {
		return
	}

	var in refundInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := in.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdBy, _ := middleware.GetUserID(c)

	refund, err := h.repo.CreateRefund(c.Request.Context(), models.Refund{
		InvoiceID:   id,
		AmountCents: in.AmountCents,
		Reason:      in.Reason,
		CreatedBy:   createdBy,
	})
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusCreated, refund)
}

func (h *SubscriptionHandler) ListRefunds(c *gin.Context) {
	id, ok := parseID(c, "invoiceID")
	if !ok {
		return
	}

	refunds, err := h.repo.ListRefunds(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, refunds)
}

// DownloadInvoicePDF streams the rendered invoice as a PDF attachment.
func (h *SubscriptionHandler) DownloadInvoicePDF(c *gin.Context) {
	id, ok := parseID(c, "invoiceID")
	if !ok {
		return
	}

	document, details, err := h.invoices.Render(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, details.InvoiceNumber))
	c.Data(http.StatusOK, "application/pdf", document)
}

// SendInvoice emails the rendered invoice to the customer.
func (h *SubscriptionHandler) SendInvoice(c *gin.Context) {
	id, ok := parseID(c, "invoiceID")
	if !ok {
		return
	}

	invoice, err := h.invoices.Send(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, invoice)
}
