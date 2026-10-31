package handlers

import (
	"errors"
	"net/http"
	"strings"

	"recur_backend/internals/models"
	"recur_backend/internals/repository"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

type productInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      *bool  `json:"active"`
}

func (in *productInput) validate() error {
	in.Name = strings.TrimSpace(in.Name)
	in.Description = strings.TrimSpace(in.Description)
	if in.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// active defaults to true when the field is omitted.
func (in productInput) active() bool {
	return in.Active == nil || *in.Active
}

func (in *productInput) toModel() models.Product {
	return models.Product{
		Name:        in.Name,
		Description: in.Description,
		Active:      in.active(),
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var in productInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := in.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.repo.CreateProduct(c.Request.Context(), in.toModel())
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.repo.ListProducts(c.Request.Context())
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	product, err := h.repo.GetProduct(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var in productInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := in.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := in.toModel()
	product.ID = id
	updated, err := h.repo.UpdateProduct(c.Request.Context(), product)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	deleted, err := h.repo.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
