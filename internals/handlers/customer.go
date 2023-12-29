package handlers

import (
	"aspire-lite/internals/constants"
	"aspire-lite/internals/models"
	"aspire-lite/internals/usecases"
	"context"
	"encoding/json"
	"net/http"
)

type CustomerRepository interface {
	GetUserCredential(ctx context.Context, username string) (*models.Customer, error)
}

type TokenEncoder interface {
	Encode(customerID int64) (string, error)
}

type authenticatorHandler struct {
	customerRepo CustomerRepository
	tokenEncoder TokenEncoder
}

func NewAuthenticator(customerRepo CustomerRepository, tokenEncoder TokenEncoder) *authenticatorHandler {
	return &authenticatorHandler{
		customerRepo: customerRepo,
		tokenEncoder: tokenEncoder,
	}
}

func (h *authenticatorHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()
	var user usecases.LoginInfo
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}

	current, err := h.customerRepo.GetUserCredential(ctx, user.Username)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if current.Password != user.Password {
		writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	token, err := h.tokenEncoder.Encode(current.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	writeOKResponse(w, map[string]interface{}{
		"token": token,
	})
}
