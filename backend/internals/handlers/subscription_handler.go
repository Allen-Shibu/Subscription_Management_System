package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"recur_backend/internals/middleware"
	"recur_backend/internals/models"
	"recur_backend/internals/repository"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionHandler(repo *repository.SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{repo: repo}
}

var subscriptionStatuses = map[string]bool{
	models.SubscriptionStatusActive:   true,
	models.SubscriptionStatusTrialing: true,
	models.SubscriptionStatusPastDue:  true,
	models.SubscriptionStatusCanceled: true,
	models.SubscriptionStatusExpired:  true,
}

type subscriptionCreateInput struct {
	CustomerID int64      `json:"customerId"`
	VariantID  int64      `json:"variantId"`
	Quantity   int        `json:"quantity"`
	StartAt    *time.Time `json:"startAt"`
}

func (in *subscriptionCreateInput) validate() error {
	if in.CustomerID <= 0 {
		return errors.New("customerId is required")
	}
	if in.VariantID <= 0 {
		return errors.New("variantId is required")
	}
	if in.Quantity == 0 {
		in.Quantity = 1
	}
	if in.Quantity < 1 {
		return errors.New("quantity must be greater than zero")
	}
	return nil
}

type subscriptionUpdateInput struct {
	Quantity           *int       `json:"quantity"`
	CurrentPeriodStart *time.Time `json:"currentPeriodStart"`
	CurrentPeriodEnd   *time.Time `json:"currentPeriodEnd"`
	CancelAtPeriodEnd  *bool      `json:"cancelAtPeriodEnd"`
}

type cancelInput struct {
	AtPeriodEnd *bool `json:"atPeriodEnd"`
}

func (in *cancelInput) validate() error {
	if in.AtPeriodEnd == nil {
		return errors.New("atPeriodEnd is required")
	}
	return nil
}

type discountInput struct {
	AmountCents int64  `json:"amountCents"`
	Reason      string `json:"reason"`
}

func (in *discountInput) validate() error {
	in.Reason = strings.TrimSpace(in.Reason)
	if in.AmountCents <= 0 {
		return errors.New("amountCents must be greater than zero")
	}
	if in.Reason == "" {
		return errors.New("reason is required")
	}
	return nil
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var in subscriptionCreateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := in.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start := time.Now().UTC()
	if in.StartAt != nil {
		start = in.StartAt.UTC()
	}
	createdBy, _ := middleware.GetUserID(c)

	subscription, err := h.repo.CreateSubscription(c.Request.Context(), models.Subscription{
		CustomerID:         in.CustomerID,
		VariantID:          in.VariantID,
		Quantity:           in.Quantity,
		CurrentPeriodStart: start,
		CreatedBy:          createdBy,
	})
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusCreated, subscription)
}

func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	var customerID int64
	if raw := c.Query("customerId"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customerId"})
			return
		}
		customerID = id
	}

	status := c.Query("status")
	if status != "" && !subscriptionStatuses[status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	subscriptions, err := h.repo.ListSubscriptions(c.Request.Context(), customerID, status)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, subscriptions)
}

func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	subscription, err := h.repo.GetSubscription(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, subscription)
}

func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var in subscriptionUpdateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	subscription, err := h.repo.GetSubscription(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}
	if in.Quantity != nil {
		if *in.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be greater than zero"})
			return
		}
		subscription.Quantity = *in.Quantity
	}
	if in.CurrentPeriodStart != nil {
		subscription.CurrentPeriodStart = in.CurrentPeriodStart.UTC()
	}
	if in.CurrentPeriodEnd != nil {
		subscription.CurrentPeriodEnd = in.CurrentPeriodEnd.UTC()
	}
	if in.CancelAtPeriodEnd != nil {
		subscription.CancelAtPeriodEnd = *in.CancelAtPeriodEnd
	}
	if !subscription.CurrentPeriodEnd.After(subscription.CurrentPeriodStart) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "currentPeriodEnd must be after currentPeriodStart"})
		return
	}

	updated, err := h.repo.UpdateSubscription(c.Request.Context(), subscription)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}

type taxRuleAssignmentInput struct {
	TaxRuleID *int64 `json:"taxRuleId"`
}

// SetSubscriptionTaxRule is admin-only and clears the rule when taxRuleId is null.
func (h *SubscriptionHandler) SetSubscriptionTaxRule(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var in taxRuleAssignmentInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if in.TaxRuleID != nil && *in.TaxRuleID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "taxRuleId must be a positive id or null"})
		return
	}

	subscription, err := h.repo.SetSubscriptionTaxRule(c.Request.Context(), id, in.TaxRuleID)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, subscription)
}

func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var in cancelInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := in.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.repo.CancelSubscription(c.Request.Context(), id, *in.AtPeriodEnd)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, subscription)
}

func (h *SubscriptionHandler) CreateDiscount(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var in discountInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := in.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdBy, _ := middleware.GetUserID(c)

	discount, err := h.repo.CreateDiscount(c.Request.Context(), models.SubscriptionDiscount{
		SubscriptionID: id,
		AmountCents:    in.AmountCents,
		Reason:         in.Reason,
		CreatedBy:      createdBy,
	})
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusCreated, discount)
}

func (h *SubscriptionHandler) ListDiscounts(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	discounts, err := h.repo.ListDiscounts(c.Request.Context(), id)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, discounts)
}
