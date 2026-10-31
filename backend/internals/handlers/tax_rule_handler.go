package handlers

import (
	"errors"
	"net/http"
	"strings"

	"recur_backend/internals/middleware"
	"recur_backend/internals/models"
	"recur_backend/internals/repository"

	"github.com/gin-gonic/gin"
)

type TaxRuleHandler struct {
	repo *repository.TaxRuleRepository
}

func NewTaxRuleHandler(repo *repository.TaxRuleRepository) *TaxRuleHandler {
	return &TaxRuleHandler{repo: repo}
}

type taxRuleInput struct {
	Name    string `json:"name"`
	RateBps int    `json:"rateBps"`
	Active  *bool  `json:"active"`
}

func (in *taxRuleInput) validate() error {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return errors.New("name is required")
	}
	if in.RateBps < 0 || in.RateBps > 10000 {
		return errors.New("rateBps must be between 0 and 10000")
	}
	return nil
}

func (in taxRuleInput) active() bool {
	return in.Active == nil || *in.Active
}

func (h *TaxRuleHandler) CreateTaxRule(c *gin.Context) {
	var in taxRuleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := in.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdBy, _ := middleware.GetUserID(c)

	rule, err := h.repo.CreateTaxRule(c.Request.Context(), models.TaxRule{
		Name:      in.Name,
		RateBps:   in.RateBps,
		Active:    in.active(),
		CreatedBy: createdBy,
	})
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusCreated, rule)
}

func (h *TaxRuleHandler) ListTaxRules(c *gin.Context) {
	rules, err := h.repo.ListTaxRules(c.Request.Context())
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, rules)
}

func (h *TaxRuleHandler) GetTaxRule(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	rule, err := h.repo.GetTaxRule(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, rule)
}

func (h *TaxRuleHandler) UpdateTaxRule(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var in taxRuleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := in.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.repo.UpdateTaxRule(c.Request.Context(), models.TaxRule{
		ID:      id,
		Name:    in.Name,
		RateBps: in.RateBps,
		Active:  in.active(),
	})
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, rule)
}

func (h *TaxRuleHandler) DeleteTaxRule(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	deleted, err := h.repo.DeleteTaxRule(c.Request.Context(), id)
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
