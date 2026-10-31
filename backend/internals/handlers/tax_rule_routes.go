package handlers

import "github.com/gin-gonic/gin"

func (h *TaxRuleHandler) RegisterRoutes(api, admin *gin.RouterGroup) {
	api.GET("/tax-rules", h.ListTaxRules)
	api.GET("/tax-rules/:id", h.GetTaxRule)

	admin.POST("/tax-rules", h.CreateTaxRule)
	admin.PUT("/tax-rules/:id", h.UpdateTaxRule)
	admin.DELETE("/tax-rules/:id", h.DeleteTaxRule)
}
