package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts all subscription endpoints on the staff/admin group.
func (h *SubscriptionHandler) RegisterRoutes(staff *gin.RouterGroup) {
	staff.POST("/subscriptions", h.CreateSubscription)
	staff.GET("/subscriptions", h.ListSubscriptions)
	staff.GET("/subscriptions/:id", h.GetSubscription)
	staff.PATCH("/subscriptions/:id", h.UpdateSubscription)
	staff.POST("/subscriptions/:id/cancel", h.CancelSubscription)
	staff.POST("/subscriptions/:id/discounts", h.CreateDiscount)
	staff.GET("/subscriptions/:id/discounts", h.ListDiscounts)
	staff.POST("/subscriptions/:id/invoices", h.CreateInvoice)
	staff.GET("/subscriptions/:id/invoices", h.ListInvoices)
	staff.GET("/invoices/:invoiceID", h.GetInvoice)
	staff.POST("/invoices/:invoiceID/pay", h.PayInvoice)
	staff.POST("/invoices/:invoiceID/refunds", h.CreateRefund)
	staff.GET("/invoices/:invoiceID/refunds", h.ListRefunds)
}
