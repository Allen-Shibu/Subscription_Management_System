package models

// DashboardSummary holds top-level KPIs for the admin dashboard.
type DashboardSummary struct {
	MRR                       int64   `json:"mrr"`
	PreviousMRR               int64   `json:"previousMrr"`
	MRRChange                 int64   `json:"mrrChange"`
	MRRChangePercent          float64 `json:"mrrChangePercent"`
	ARR                       int64   `json:"arr"`
	ActiveSubscriptions       int64   `json:"activeSubscriptions"`
	TrialingSubscriptions     int64   `json:"trialingSubscriptions"`
	PastDueSubscriptions      int64   `json:"pastDueSubscriptions"`
	CanceledSubscriptions     int64   `json:"canceledSubscriptions"`
	NewSubscriptionsThisMonth int64   `json:"newSubscriptionsThisMonth"`
	CanceledThisMonth         int64   `json:"canceledThisMonth"`
	TotalPaidRevenue          int64   `json:"totalPaidRevenue"`
	OutstandingRevenue        int64   `json:"outstandingRevenue"`
	TotalRefunded             int64   `json:"totalRefunded"`
	Currency                  string  `json:"currency"`
}

// MonthlyPoint is one bucket of the time series, usable for line/bar charts.
type MonthlyPoint struct {
	Month        string `json:"month"` // YYYY-MM
	MRR          int64  `json:"mrr"`
	Revenue      int64  `json:"revenue"` // invoiced total
	InvoiceCount int64  `json:"invoiceCount"`
}

// PlanBreakdown is MRR grouped by product variant, usable for pie/bar charts.
type PlanBreakdown struct {
	VariantID         int64  `json:"variantId"`
	VariantName       string `json:"variantName"`
	ProductName       string `json:"productName"`
	MRR               int64  `json:"mrr"`
	SubscriptionCount int64  `json:"subscriptionCount"`
}

// DashboardResponse is the full payload for the admin dashboard.
type DashboardResponse struct {
	Summary DashboardSummary `json:"summary"`
	Monthly []MonthlyPoint   `json:"monthly"`
	ByPlan  []PlanBreakdown  `json:"byPlan"`
}
