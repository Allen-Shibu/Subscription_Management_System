package repository

import (
	"context"
	"time"

	"recur_backend/internals/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetDashboard(ctx context.Context, months int) (models.DashboardResponse, error) {
	if months < 1 {
		months = 12
	}
	if months > 24 {
		months = 24
	}

	resp := models.DashboardResponse{
		Monthly: []models.MonthlyPoint{},
		ByPlan:  []models.PlanBreakdown{},
	}

	// Monthly MRR time series.
	mrrRows, err := r.db.Query(ctx, `
		WITH months AS (
			SELECT generate_series(
				date_trunc('month', NOW() - ($1 - 1) * INTERVAL '1 month'),
				date_trunc('month', NOW()),
				INTERVAL '1 month'
			) AS month_start
		)
		SELECT to_char(m.month_start, 'YYYY-MM') AS month,
			COALESCE(SUM(CASE WHEN s.billing_interval = 'year'
				THEN s.unit_price_cents * s.quantity / 12
				ELSE s.unit_price_cents * s.quantity
			END), 0) AS mrr
		FROM months m
		LEFT JOIN subscriptions s
			ON s.created_at < (m.month_start + INTERVAL '1 month')
			AND (s.ended_at IS NULL OR s.ended_at > (m.month_start + INTERVAL '1 month'))
			AND (s.canceled_at IS NULL OR s.canceled_at > (m.month_start + INTERVAL '1 month'))
			AND s.status IN ('active', 'trialing', 'past_due')
		GROUP BY m.month_start
		ORDER BY m.month_start
	`, months)
	if err != nil {
		return resp, err
	}
	mrrByMonth := map[string]int64{}
	for mrrRows.Next() {
		var month string
		var mrr int64
		if err := mrrRows.Scan(&month, &mrr); err != nil {
			mrrRows.Close()
			return resp, err
		}
		mrrByMonth[month] = mrr
	}
	mrrRows.Close()
	if err := mrrRows.Err(); err != nil {
		return resp, err
	}

	// Monthly invoiced revenue.
	revRows, err := r.db.Query(ctx, `
		SELECT to_char(date_trunc('month', issued_at), 'YYYY-MM') AS month,
			COALESCE(SUM(total_cents), 0) AS revenue,
			COUNT(*) AS invoice_count
		FROM invoices
		WHERE issued_at >= date_trunc('month', NOW() - ($1 - 1) * INTERVAL '1 month')
		GROUP BY 1
	`, months)
	if err != nil {
		return resp, err
	}
	revenueByMonth := map[string]struct {
		revenue int64
		count   int64
	}{}
	for revRows.Next() {
		var month string
		var revenue, count int64
		if err := revRows.Scan(&month, &revenue, &count); err != nil {
			revRows.Close()
			return resp, err
		}
		revenueByMonth[month] = struct {
			revenue int64
			count   int64
		}{revenue, count}
	}
	revRows.Close()
	if err := revRows.Err(); err != nil {
		return resp, err
	}

	// Build ascending month list.
	now := time.Now().UTC()
	for i := months - 1; i >= 0; i-- {
		t := now.AddDate(0, -i, 0)
		month := t.Format("2006-01")
		point := models.MonthlyPoint{Month: month, MRR: mrrByMonth[month]}
		if v, ok := revenueByMonth[month]; ok {
			point.Revenue = v.revenue
			point.InvoiceCount = v.count
		}
		resp.Monthly = append(resp.Monthly, point)
	}

	// Summary counts.
	var active, trialing, pastDue, canceled int64
	err = r.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'active'),
			COUNT(*) FILTER (WHERE status = 'trialing'),
			COUNT(*) FILTER (WHERE status = 'past_due'),
			COUNT(*) FILTER (WHERE status = 'canceled')
		FROM subscriptions`).Scan(&active, &trialing, &pastDue, &canceled)
	if err != nil {
		return resp, err
	}

	var newThisMonth, canceledThisMonth int64
	err = r.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE created_at >= date_trunc('month', NOW())),
			COUNT(*) FILTER (WHERE canceled_at >= date_trunc('month', NOW()))
		FROM subscriptions`).Scan(&newThisMonth, &canceledThisMonth)
	if err != nil {
		return resp, err
	}

	var totalPaid, outstanding, totalRefunded int64
	err = r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(total_cents) FILTER (WHERE status = 'paid'), 0),
			COALESCE(SUM(total_cents) FILTER (WHERE status = 'open'), 0)
		FROM invoices`).Scan(&totalPaid, &outstanding)
	if err != nil {
		return resp, err
	}
	err = r.db.QueryRow(ctx, `SELECT COALESCE(SUM(amount_cents), 0) FROM refunds`).Scan(&totalRefunded)
	if err != nil {
		return resp, err
	}

	// Currency: use the single currency if all active subs share one.
	var currencyCount int
	var currency string
	err = r.db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT currency), COALESCE(MIN(currency), 'USD')
		FROM subscriptions WHERE status IN ('active','trialing','past_due') AND ended_at IS NULL`).Scan(&currencyCount, &currency)
	if err != nil {
		return resp, err
	}
	if currencyCount != 1 {
		currency = "mixed"
	}
	if currency == "" {
		currency = "USD"
	}

	currentMRR := int64(0)
	previousMRR := int64(0)
	if len(resp.Monthly) > 0 {
		currentMRR = resp.Monthly[len(resp.Monthly)-1].MRR
	}
	if len(resp.Monthly) > 1 {
		previousMRR = resp.Monthly[len(resp.Monthly)-2].MRR
	}
	change := currentMRR - previousMRR
	changePercent := 0.0
	if previousMRR > 0 {
		changePercent = float64(change) / float64(previousMRR) * 100
	}

	resp.Summary = models.DashboardSummary{
		MRR:                       currentMRR,
		PreviousMRR:               previousMRR,
		MRRChange:                 change,
		MRRChangePercent:          changePercent,
		ARR:                       currentMRR * 12,
		ActiveSubscriptions:       active,
		TrialingSubscriptions:     trialing,
		PastDueSubscriptions:      pastDue,
		CanceledSubscriptions:     canceled,
		NewSubscriptionsThisMonth: newThisMonth,
		CanceledThisMonth:         canceledThisMonth,
		TotalPaidRevenue:          totalPaid,
		OutstandingRevenue:        outstanding,
		TotalRefunded:             totalRefunded,
		Currency:                  currency,
	}

	// By plan.
	planRows, err := r.db.Query(ctx, `
		SELECT s.variant_id, pv.name, p.name,
			COALESCE(SUM(CASE WHEN s.billing_interval = 'year'
				THEN s.unit_price_cents * s.quantity / 12
				ELSE s.unit_price_cents * s.quantity
			END), 0) AS mrr,
			COUNT(*) AS subscription_count
		FROM subscriptions s
		JOIN product_variants pv ON pv.id = s.variant_id
		JOIN products p ON p.id = pv.product_id
		WHERE s.status IN ('active', 'trialing', 'past_due') AND s.ended_at IS NULL
		GROUP BY s.variant_id, pv.name, p.name
		ORDER BY mrr DESC
	`)
	if err != nil {
		return resp, err
	}
	for planRows.Next() {
		var pb models.PlanBreakdown
		if err := planRows.Scan(&pb.VariantID, &pb.VariantName, &pb.ProductName, &pb.MRR, &pb.SubscriptionCount); err != nil {
			planRows.Close()
			return resp, err
		}
		resp.ByPlan = append(resp.ByPlan, pb)
	}
	planRows.Close()
	if err := planRows.Err(); err != nil {
		return resp, err
	}

	return resp, nil
}
