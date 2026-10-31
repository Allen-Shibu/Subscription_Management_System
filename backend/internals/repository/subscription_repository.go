package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"recur_backend/internals/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvoiceNotPaid     = errors.New("invoice is not paid")
	ErrRefundExceedsTotal = errors.New("refund exceeds invoice balance")
)

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

const subscriptionColumns = `id, customer_id, variant_id, status, quantity, unit_price_cents, currency,
	billing_interval, current_period_start, current_period_end, cancel_at_period_end,
	canceled_at, ended_at, created_by, created_at, updated_at, tax_rule_id`

const invoiceColumns = `id, invoice_number, subscription_id, status, period_start, period_end,
	subtotal_cents, discount_cents, total_cents, issued_at, due_at, paid_at, email_sent_at,
	created_by, created_at, tax_rule_id, tax_rate_bps, tax_cents`

const discountColumns = `id, subscription_id, amount_cents, reason, applied_invoice_id, created_by, created_at`

const refundColumns = `id, invoice_id, amount_cents, reason, refunded_at, created_by, created_at`

func scanSubscription(row pgx.Row) (models.Subscription, error) {
	var s models.Subscription
	err := row.Scan(&s.ID, &s.CustomerID, &s.VariantID, &s.Status, &s.Quantity,
		&s.UnitPriceCents, &s.Currency, &s.BillingInterval, &s.CurrentPeriodStart,
		&s.CurrentPeriodEnd, &s.CancelAtPeriodEnd, &s.CanceledAt, &s.EndedAt,
		&s.CreatedBy, &s.CreatedAt, &s.UpdatedAt, &s.TaxRuleID)
	return s, err
}

func scanInvoice(row pgx.Row) (models.Invoice, error) {
	var i models.Invoice
	err := row.Scan(&i.ID, &i.InvoiceNumber, &i.SubscriptionID, &i.Status, &i.PeriodStart,
		&i.PeriodEnd, &i.SubtotalCents, &i.DiscountCents, &i.TotalCents, &i.IssuedAt,
		&i.DueAt, &i.PaidAt, &i.EmailSentAt, &i.CreatedBy, &i.CreatedAt,
		&i.TaxRuleID, &i.TaxRateBps, &i.TaxCents)
	return i, err
}

func scanDiscount(row pgx.Row) (models.SubscriptionDiscount, error) {
	var d models.SubscriptionDiscount
	err := row.Scan(&d.ID, &d.SubscriptionID, &d.AmountCents, &d.Reason,
		&d.AppliedInvoiceID, &d.CreatedBy, &d.CreatedAt)
	return d, err
}

func scanRefund(row pgx.Row) (models.Refund, error) {
	var r models.Refund
	err := row.Scan(&r.ID, &r.InvoiceID, &r.AmountCents, &r.Reason,
		&r.RefundedAt, &r.CreatedBy, &r.CreatedAt)
	return r, err
}

// CreateSubscription snapshots price/currency/interval from the variant and
// computes the first period end from that interval.
func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, s models.Subscription) (models.Subscription, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO subscriptions (customer_id, variant_id, status, quantity, unit_price_cents,
			currency, billing_interval, current_period_start, current_period_end, created_by)
		SELECT $1, v.id, 'active', $2, v.price_cents, v.currency, v.billing_interval,
			$3, $3 + CASE v.billing_interval
				WHEN 'month' THEN INTERVAL '1 month'
				ELSE INTERVAL '1 year'
			END,
			$4
		FROM product_variants v
		JOIN products p ON p.id = v.product_id
		WHERE v.id = $5 AND v.is_active AND p.is_active
		RETURNING `+subscriptionColumns,
		s.CustomerID, s.Quantity, s.CurrentPeriodStart, s.CreatedBy, s.VariantID)
	return scanSubscription(row)
}

func (r *SubscriptionRepository) ListSubscriptions(ctx context.Context, customerID int64, status string) ([]models.Subscription, error) {
	query := `SELECT ` + subscriptionColumns + ` FROM subscriptions`
	args := []any{}
	clauses := []string{}

	if customerID > 0 {
		args = append(args, customerID)
		clauses = append(clauses, fmt.Sprintf("customer_id = $%d", len(args)))
	}
	if status != "" {
		args = append(args, status)
		clauses = append(clauses, fmt.Sprintf("status = $%d", len(args)))
	}
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	query += " ORDER BY id"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscriptions := []models.Subscription{}
	for rows.Next() {
		s, err := scanSubscription(rows)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}
	return subscriptions, rows.Err()
}

func (r *SubscriptionRepository) GetSubscription(ctx context.Context, id int64) (models.Subscription, error) {
	row := r.db.QueryRow(ctx, `SELECT `+subscriptionColumns+` FROM subscriptions WHERE id = $1`, id)
	return scanSubscription(row)
}

func (r *SubscriptionRepository) UpdateSubscription(ctx context.Context, s models.Subscription) (models.Subscription, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE subscriptions
		SET quantity = $2, current_period_start = $3, current_period_end = $4,
			cancel_at_period_end = $5,
			canceled_at = CASE WHEN $5 THEN COALESCE(canceled_at, NOW()) ELSE NULL END,
			updated_at = NOW()
		WHERE id = $1
		RETURNING `+subscriptionColumns,
		s.ID, s.Quantity, s.CurrentPeriodStart, s.CurrentPeriodEnd, s.CancelAtPeriodEnd)
	return scanSubscription(row)
}

// SetSubscriptionTaxRule attaches or clears the admin-defined tax rule.
func (r *SubscriptionRepository) SetSubscriptionTaxRule(ctx context.Context, id int64, taxRuleID *int64) (models.Subscription, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE subscriptions
		SET tax_rule_id = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING `+subscriptionColumns, id, taxRuleID)
	return scanSubscription(row)
}

// CancelSubscription either flags cancellation at period end or cancels now.
func (r *SubscriptionRepository) CancelSubscription(ctx context.Context, id int64, atPeriodEnd bool) (models.Subscription, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE subscriptions
		SET cancel_at_period_end = $2,
			canceled_at = CASE WHEN $2 THEN COALESCE(canceled_at, NOW()) ELSE NOW() END,
			ended_at = CASE WHEN $2 THEN NULL ELSE NOW() END,
			status = CASE WHEN $2 THEN status ELSE 'canceled' END,
			updated_at = NOW()
		WHERE id = $1
		RETURNING `+subscriptionColumns, id, atPeriodEnd)
	return scanSubscription(row)
}

func (r *SubscriptionRepository) CreateDiscount(ctx context.Context, d models.SubscriptionDiscount) (models.SubscriptionDiscount, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO subscription_discounts (subscription_id, amount_cents, reason, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING `+discountColumns,
		d.SubscriptionID, d.AmountCents, d.Reason, d.CreatedBy)
	return scanDiscount(row)
}

func (r *SubscriptionRepository) ListDiscounts(ctx context.Context, subscriptionID int64) ([]models.SubscriptionDiscount, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+discountColumns+` FROM subscription_discounts WHERE subscription_id = $1 ORDER BY id`,
		subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	discounts := []models.SubscriptionDiscount{}
	for rows.Next() {
		d, err := scanDiscount(rows)
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, d)
	}
	return discounts, rows.Err()
}

// CreateInvoice issues one invoice for the subscription's current period and
// consumes all pending one-time discounts.
func (r *SubscriptionRepository) CreateInvoice(ctx context.Context, subscriptionID int64, createdBy string) (models.Invoice, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return models.Invoice{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sub, err := scanSubscription(tx.QueryRow(ctx,
		`SELECT `+subscriptionColumns+` FROM subscriptions WHERE id = $1 FOR UPDATE`, subscriptionID))
	if err != nil {
		return models.Invoice{}, err
	}

	subtotal := sub.UnitPriceCents * int64(sub.Quantity)

	var discount int64
	err = tx.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount_cents), 0)
		FROM subscription_discounts
		WHERE subscription_id = $1 AND applied_invoice_id IS NULL`, subscriptionID).Scan(&discount)
	if err != nil {
		return models.Invoice{}, err
	}
	if discount > subtotal {
		discount = subtotal
	}

	taxable := subtotal - discount
	var taxRuleID *int64
	var taxRate int
	var taxCents int64
	if sub.TaxRuleID != nil {
		if err := tx.QueryRow(ctx,
			`SELECT rate_bps FROM tax_rules WHERE id = $1`, *sub.TaxRuleID).Scan(&taxRate); err != nil {
			return models.Invoice{}, err
		}
		taxRuleID = sub.TaxRuleID
		taxCents = taxCentsFor(taxable, taxRate)
	}
	total := taxable + taxCents

	invoice, err := scanInvoice(tx.QueryRow(ctx, `
		INSERT INTO invoices (subscription_id, status, period_start, period_end, subtotal_cents,
			discount_cents, tax_rule_id, tax_rate_bps, tax_cents, total_cents, issued_at, due_at, created_by)
		VALUES ($1, 'open', $2, $3, $4, $5, $6, $7, $8, $9, NOW(), $3, $10)
		RETURNING `+invoiceColumns,
		subscriptionID, sub.CurrentPeriodStart, sub.CurrentPeriodEnd,
		subtotal, discount, taxRuleID, taxRate, taxCents, total, createdBy))
	if err != nil {
		return models.Invoice{}, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE subscription_discounts
		SET applied_invoice_id = $1
		WHERE subscription_id = $2 AND applied_invoice_id IS NULL`, invoice.ID, subscriptionID)
	if err != nil {
		return models.Invoice{}, err
	}

	// Move the subscription to the next billing period so a later invoice
	// call bills the following period.
	_, err = tx.Exec(ctx, `
		UPDATE subscriptions
		SET current_period_start = current_period_end,
			current_period_end = current_period_end + CASE billing_interval
				WHEN 'month' THEN INTERVAL '1 month'
				ELSE INTERVAL '1 year'
			END,
			updated_at = NOW()
		WHERE id = $1`, subscriptionID)
	if err != nil {
		return models.Invoice{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Invoice{}, err
	}
	return invoice, nil
}

func (r *SubscriptionRepository) ListInvoices(ctx context.Context, subscriptionID int64) ([]models.Invoice, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+invoiceColumns+` FROM invoices WHERE subscription_id = $1 ORDER BY id`, subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := []models.Invoice{}
	for rows.Next() {
		i, err := scanInvoice(rows)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, i)
	}
	return invoices, rows.Err()
}

func (r *SubscriptionRepository) GetInvoice(ctx context.Context, id int64) (models.Invoice, error) {
	row := r.db.QueryRow(ctx, `SELECT `+invoiceColumns+` FROM invoices WHERE id = $1`, id)
	return scanInvoice(row)
}

// GetInvoiceDetails loads everything needed to render an invoice document.
func (r *SubscriptionRepository) GetInvoiceDetails(ctx context.Context, id int64) (models.InvoiceDetails, error) {
	var d models.InvoiceDetails
	err := r.db.QueryRow(ctx, `
		SELECT i.id, i.invoice_number, i.subscription_id, i.status, i.period_start, i.period_end,
			i.subtotal_cents, i.discount_cents, i.total_cents, i.issued_at, i.due_at, i.paid_at,
			i.email_sent_at, i.created_by, i.created_at,
			COALESCE(NULLIF(BTRIM(CONCAT_WS(' ', u.first_name, u.last_name)), ''), u.email, 'Customer'),
			COALESCE(u.email, ''),
			p.name, v.name, s.quantity, s.unit_price_cents, s.currency, s.billing_interval
		FROM invoices i
		JOIN subscriptions s ON s.id = i.subscription_id
		JOIN users u ON u.id = s.customer_id
		JOIN product_variants v ON v.id = s.variant_id
		JOIN products p ON p.id = v.product_id
		WHERE i.id = $1`, id).
		Scan(&d.ID, &d.InvoiceNumber, &d.SubscriptionID, &d.Status, &d.PeriodStart, &d.PeriodEnd,
			&d.SubtotalCents, &d.DiscountCents, &d.TotalCents, &d.IssuedAt, &d.DueAt, &d.PaidAt,
			&d.EmailSentAt, &d.CreatedBy, &d.CreatedAt,
			&d.CustomerName, &d.CustomerEmail, &d.ProductName, &d.VariantName,
			&d.Quantity, &d.UnitPriceCents, &d.Currency, &d.BillingInterval)
	return d, err
}

// MarkInvoiceEmailSent records a successful delivery.
func (r *SubscriptionRepository) MarkInvoiceEmailSent(ctx context.Context, id int64) (models.Invoice, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE invoices SET email_sent_at = NOW()
		WHERE id = $1
		RETURNING `+invoiceColumns, id)
	return scanInvoice(row)
}

func (r *SubscriptionRepository) MarkInvoicePaid(ctx context.Context, id int64) (models.Invoice, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE invoices
		SET status = 'paid', paid_at = NOW()
		WHERE id = $1 AND status = 'open'
		RETURNING `+invoiceColumns, id)
	return scanInvoice(row)
}

// CreateRefund records a refund against a paid invoice, refusing to refund
// more than the invoice total.
func (r *SubscriptionRepository) CreateRefund(ctx context.Context, rf models.Refund) (models.Refund, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return models.Refund{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var total int64
	var status string
	err = tx.QueryRow(ctx,
		`SELECT total_cents, status FROM invoices WHERE id = $1 FOR UPDATE`, rf.InvoiceID).
		Scan(&total, &status)
	if err != nil {
		return models.Refund{}, err
	}
	if status != models.InvoiceStatusPaid {
		return models.Refund{}, ErrInvoiceNotPaid
	}

	var refunded int64
	err = tx.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount_cents), 0) FROM refunds WHERE invoice_id = $1`, rf.InvoiceID).
		Scan(&refunded)
	if err != nil {
		return models.Refund{}, err
	}
	if rf.AmountCents > total-refunded {
		return models.Refund{}, ErrRefundExceedsTotal
	}

	row := tx.QueryRow(ctx, `
		INSERT INTO refunds (invoice_id, amount_cents, reason, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING `+refundColumns,
		rf.InvoiceID, rf.AmountCents, rf.Reason, rf.CreatedBy)

	refund, err := scanRefund(row)
	if err != nil {
		return models.Refund{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Refund{}, err
	}
	return refund, nil
}

func (r *SubscriptionRepository) ListRefunds(ctx context.Context, invoiceID int64) ([]models.Refund, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+refundColumns+` FROM refunds WHERE invoice_id = $1 ORDER BY id`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	refunds := []models.Refund{}
	for rows.Next() {
		rf, err := scanRefund(rows)
		if err != nil {
			return nil, err
		}
		refunds = append(refunds, rf)
	}
	return refunds, rows.Err()
}

// taxCentsFor rounds half up.
func taxCentsFor(taxable int64, rateBps int) int64 {
	if taxable <= 0 || rateBps <= 0 {
		return 0
	}
	return (taxable*int64(rateBps) + 5000) / 10000
}
