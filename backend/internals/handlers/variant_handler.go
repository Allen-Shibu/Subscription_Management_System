package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"recur_backend/internals/models"

	"github.com/gin-gonic/gin"
)

var currencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)

var billingIntervals = map[string]bool{
	"month": true,
	"year":  true,
}

type variantInput struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	PriceCents      int64  `json:"priceCents"`
	Currency        string `json:"currency"`
	BillingInterval string `json:"billingInterval"`
	Active          *bool  `json:"active"`
}

func (in *variantInput) validate() error {
	in.Name = strings.TrimSpace(in.Name)
	in.Description = strings.TrimSpace(in.Description)
	in.BillingInterval = strings.ToLower(strings.TrimSpace(in.BillingInterval))
	in.Currency = strings.ToUpper(strings.TrimSpace(in.Currency))

	if in.Name == "" {
		return errors.New("name is required")
	}
	if in.PriceCents < 0 {
		return errors.New("priceCents must be zero or greater")
	}
	if !billingIntervals[in.BillingInterval] {
		return errors.New("billingInterval must be month or year")
	}
	if in.Currency == "" {
		in.Currency = "USD"
	}
	if !currencyPattern.MatchString(in.Currency) {
		return errors.New("currency must be a 3-letter code")
	}
	return nil
}

// active defaults to true when the field is omitted.
func (in variantInput) active() bool {
	return in.Active == nil || *in.Active
}

func (in *variantInput) toModel() models.ProductVariant {
	return models.ProductVariant{
		Name:            in.Name,
		Description:     in.Description,
		PriceCents:      in.PriceCents,
		Currency:        in.Currency,
		BillingInterval: in.BillingInterval,
		Active:          in.active(),
	}
}

func (h *ProductHandler) CreateVariant(c *gin.Context) {
	productID, ok := parseID(c, "id")
	if !ok {
		return
	}

	var in variantInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := in.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variant := in.toModel()
	variant.ProductID = productID
	created, err := h.repo.CreateVariant(c.Request.Context(), variant)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *ProductHandler) ListVariants(c *gin.Context) {
	productID, ok := parseID(c, "id")
	if !ok {
		return
	}

	variants, err := h.repo.ListVariants(c.Request.Context(), productID)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, variants)
}

func (h *ProductHandler) GetVariant(c *gin.Context) {
	id, ok := parseID(c, "variantID")
	if !ok {
		return
	}

	variant, err := h.repo.GetVariant(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, variant)
}

func (h *ProductHandler) UpdateVariant(c *gin.Context) {
	id, ok := parseID(c, "variantID")
	if !ok {
		return
	}

	var in variantInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := in.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variant := in.toModel()
	variant.ID = id
	updated, err := h.repo.UpdateVariant(c.Request.Context(), variant)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *ProductHandler) DeleteVariant(c *gin.Context) {
	id, ok := parseID(c, "variantID")
	if !ok {
		return
	}

	deleted, err := h.repo.DeleteVariant(c.Request.Context(), id)
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
