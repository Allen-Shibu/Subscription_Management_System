package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts read endpoints on the authenticated group and
// write endpoints on the admin-only group.
func (h *ProductHandler) RegisterRoutes(api, admin *gin.RouterGroup) {
	api.GET("/products", h.ListProducts)
	api.GET("/products/:id", h.GetProduct)
	api.GET("/products/:id/variants", h.ListVariants)
	api.GET("/variants/:variantID", h.GetVariant)

	admin.POST("/products", h.CreateProduct)
	admin.PUT("/products/:id", h.UpdateProduct)
	admin.DELETE("/products/:id", h.DeleteProduct)
	admin.POST("/products/:id/variants", h.CreateVariant)
	admin.PUT("/variants/:variantID", h.UpdateVariant)
	admin.DELETE("/variants/:variantID", h.DeleteVariant)
}
