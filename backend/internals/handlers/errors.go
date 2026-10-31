package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"recur_backend/internals/repository"
	invoicesvc "recur_backend/internals/services/invoice"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func parseID(c *gin.Context, param string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(param), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + param})
		return 0, false
	}
	return id, true
}

// respondDBError maps PostgreSQL errors onto HTTP status codes.
func respondDBError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.Is(err, repository.ErrInvoiceNotPaid):
		c.JSON(http.StatusConflict, gin.H{"error": "invoice must be paid before refunding"})
	case errors.Is(err, repository.ErrRefundExceedsTotal):
		c.JSON(http.StatusBadRequest, gin.H{"error": "refund exceeds the remaining invoice total"})
	case errors.Is(err, invoicesvc.ErrInvoiceAlreadySent):
		c.JSON(http.StatusConflict, gin.H{"error": "invoice has already been emailed"})
	case errors.Is(err, invoicesvc.ErrInvoiceNotSendable):
		c.JSON(http.StatusConflict, gin.H{"error": "invoice cannot be emailed in its current status"})
	case errors.Is(err, invoicesvc.ErrCustomerEmailMissing):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "customer does not have an email address"})
	case errors.Is(err, invoicesvc.ErrEmailDelivery):
		log.Printf("[ERROR] email: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to send invoice email"})
	case isPGError(err, "23505"):
		c.JSON(http.StatusConflict, gin.H{"error": "already exists"})
	case isPGError(err, "23503"):
		c.JSON(http.StatusNotFound, gin.H{"error": "referenced record not found"})
	default:
		log.Printf("[ERROR] database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func isPGError(err error, code string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code
}
