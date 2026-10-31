package models

import "time"

const (
	SubscriptionStatusActive   = "active"
	SubscriptionStatusTrialing = "trialing"
	SubscriptionStatusPastDue  = "past_due"
	SubscriptionStatusCanceled = "canceled"
	SubscriptionStatusExpired  = "expired"
)

const (
	InvoiceStatusDraft         = "draft"
	InvoiceStatusOpen          = "open"
	InvoiceStatusPaid          = "paid"
	InvoiceStatusVoid          = "void"
	InvoiceStatusUncollectible = "uncollectible"
)

const (
	BillingIntervalMonth = "month"
	BillingIntervalYear  = "year"
)

// Subscription is a recurring charge for one customer on one product tier.
// Price and interval are snapshots taken when the subscription starts.
type Subscription struct {
	ID                 int64      `json:"id"`
	CustomerID         int64      `json:"customerId"`
	VariantID          int64      `json:"variantId"`
	Status             string     `json:"status"`
	Quantity           int        `json:"quantity"`
	UnitPriceCents     int64      `json:"unitPriceCents"`
	Currency           string     `json:"currency"`
	BillingInterval    string     `json:"billingInterval"` // month | year
	CurrentPeriodStart time.Time  `json:"currentPeriodStart"`
	CurrentPeriodEnd   time.Time  `json:"currentPeriodEnd"`
	CancelAtPeriodEnd  bool       `json:"cancelAtPeriodEnd"`
	CanceledAt         *time.Time `json:"canceledAt,omitempty"`
	EndedAt            *time.Time `json:"endedAt,omitempty"`
	CreatedBy          string     `json:"createdBy"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}

// Invoice is one billing period of a subscription. EmailSentAt is set once
// the invoice has been emailed to the customer.
type Invoice struct {
	ID             int64      `json:"id"`
	InvoiceNumber  string     `json:"invoiceNumber"`
	SubscriptionID int64      `json:"subscriptionId"`
	Status         string     `json:"status"`
	PeriodStart    time.Time  `json:"periodStart"`
	PeriodEnd      time.Time  `json:"periodEnd"`
	SubtotalCents  int64      `json:"subtotalCents"`
	DiscountCents  int64      `json:"discountCents"`
	TotalCents     int64      `json:"totalCents"`
	IssuedAt       *time.Time `json:"issuedAt,omitempty"`
	DueAt          *time.Time `json:"dueAt,omitempty"`
	PaidAt         *time.Time `json:"paidAt,omitempty"`
	EmailSentAt    *time.Time `json:"emailSentAt,omitempty"`
	CreatedBy      string     `json:"createdBy"`
	CreatedAt      time.Time  `json:"createdAt"`
}

// SubscriptionDiscount is a one-time discount that is consumed by the next
// invoice; AppliedInvoiceID is nil while still pending.
type SubscriptionDiscount struct {
	ID               int64     `json:"id"`
	SubscriptionID   int64     `json:"subscriptionId"`
	AmountCents      int64     `json:"amountCents"`
	Reason           string    `json:"reason"`
	AppliedInvoiceID *int64    `json:"appliedInvoiceId,omitempty"`
	CreatedBy        string    `json:"createdBy"`
	CreatedAt        time.Time `json:"createdAt"`
}

// Refund records money returned to a customer against a paid invoice.
type Refund struct {
	ID          int64     `json:"id"`
	InvoiceID   int64     `json:"invoiceId"`
	AmountCents int64     `json:"amountCents"`
	Reason      string    `json:"reason"`
	RefundedAt  time.Time `json:"refundedAt"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
}
