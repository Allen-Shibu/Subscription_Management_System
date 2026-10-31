package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"recur_backend/pkg/types"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/invitation"
	"github.com/gin-gonic/gin"
)

type staffInviteInput struct {
	Email string `json:"email"`
}

func (in *staffInviteInput) validate() error {
	in.Email = strings.TrimSpace(in.Email)
	addr, err := mail.ParseAddress(in.Email)
	if err != nil || addr.Address != in.Email {
		return errors.New("a valid email is required")
	}
	return nil
}

// CreateStaff invites an email address to Clerk with the staff role. The
// invited user sets their own password, so no credentials pass through us.
func CreateStaff(c *gin.Context) {
	var in staffInviteInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := in.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metadata, err := json.Marshal(map[string]string{"role": types.RoleStaff})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	raw := json.RawMessage(metadata)

	inv, err := invitation.Create(c.Request.Context(), &invitation.CreateParams{
		EmailAddress:   in.Email,
		PublicMetadata: &raw,
		Notify:         clerk.Bool(true),
	})
	if err != nil {
		var apiErr *clerk.APIErrorResponse
		if errors.As(err, &apiErr) {
			for _, e := range apiErr.Errors {
				if e.Code == "form_identifier_exists" {
					c.JSON(http.StatusConflict, gin.H{"error": "a user with this email already exists"})
					return
				}
			}
		}
		log.Printf("[ERROR] clerk invitation: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to invite staff member"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           inv.ID,
		"emailAddress": inv.EmailAddress,
		"status":       inv.Status,
		"url":          inv.URL,
	})
}

func RegisterStaffRoutes(admin *gin.RouterGroup) {
	admin.POST("/staff", CreateStaff)
}
