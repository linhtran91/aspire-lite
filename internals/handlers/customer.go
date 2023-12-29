package handlers

import (
	"aspire-lite/internals/models"
	"aspire-lite/internals/usecases"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CustomerRepository interface {
	GetUserCredential(ctx context.Context, username string) (*models.Customer, error)
}

type authenticatorHandler struct {
	customerRepo CustomerRepository
	secretKey    string
}

func NewAuthenticator(customerRepo CustomerRepository, secretKey string) *authenticatorHandler {
	return &authenticatorHandler{
		customerRepo: customerRepo,
		secretKey:    secretKey,
	}
}

func (h *authenticatorHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var user usecases.LoginInfo
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}

	current, err := h.customerRepo.GetUserCredential(ctx, user.Username)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error2222")
		return
	}

	if current.Password != user.Password {
		writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	token, err := generateToken(h.secretKey, current.ID)
	if err != nil {
		fmt.Println(err)
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error3333")
		return
	}
	writeOKResponse(w, map[string]interface{}{
		"token": token,
	})
}
