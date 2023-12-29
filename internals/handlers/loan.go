package handlers

import (
	"aspire-lite/internals/constants"
	"aspire-lite/internals/models"
	"aspire-lite/internals/usecases"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type LoanRepository interface {
	Create(ctx context.Context, loan *models.Loan, repayments []*models.Repayment) (int64, error)
	Approve(ctx context.Context, loanID int64, at time.Time) error
	View(ctx context.Context, customerID int64, limit, offset int) ([]*models.Loan, error)
	UpdateStatus(ctx context.Context, loanID int64, at time.Time) error
}

type TokenDecoder interface {
	Decode(tokenString string) (int64, error)
}

type loanHandler struct {
	loanRepo     LoanRepository
	tokenDecoder TokenDecoder
	customerRepo CustomerRepository
}

func NewLoan(loanRepo LoanRepository, tokenDecoder TokenDecoder, customerRepo CustomerRepository) *loanHandler {
	return &loanHandler{
		customerRepo: customerRepo,
		loanRepo:     loanRepo,
		tokenDecoder: tokenDecoder,
	}
}

func (h *loanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()
	var loan usecases.Loan
	if err := json.NewDecoder(r.Body).Decode(&loan); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}
	defer r.Body.Close()

	inputs := mux.Vars(r)
	customerID, err := strconv.Atoi(inputs["customer_id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}

	_, err = h.customerRepo.GetUserByID(ctx, int64(customerID))
	if err != nil && err == constants.ErrorRecordNotFound {
		writeErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	now := time.Now().UTC()
	date, err := parseDate(loan.Date)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}
	newLoan := &models.Loan{
		Amount:        loan.Amount,
		Term:          loan.Term,
		CustomerID:    int64(customerID),
		Status:        constants.PENDING,
		ScheduledDate: date,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	repayments := make([]*models.Repayment, loan.Term)
	var total float64
	for i := 0; i < loan.Term; i++ {
		amount := float64(int((loan.Amount/float64(loan.Term))*100)) / 100
		total += amount
		repayments[i] = &models.Repayment{
			ID:              generateUUID(),
			ScheduledAmount: amount,
			Status:          constants.PENDING,
			ScheduledPayAt:  date.AddDate(0, 0, (i+1)*7),
			CreatedAt:       now,
			UpdatedAt:       now,
		}
	}
	repayments[len(repayments)-1].ScheduledAmount += loan.Amount - total

	id, err := h.loanRepo.Create(ctx, newLoan, repayments)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	writeOKResponse(w, map[string]interface{}{
		"loan_id":    id,
		"repayments": repayments,
	})
}

func (h *loanHandler) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()
	inputs := mux.Vars(r)
	loanID, err := strconv.Atoi(inputs["loan_id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}

	if err := h.loanRepo.Approve(ctx, int64(loanID), time.Now().UTC()); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	writeOKResponse(w, map[string]interface{}{
		"loan_id": loanID,
	})
}

func (h *loanHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()

	token := r.Header.Get("Authorization")
	token = strings.ReplaceAll(token, "Bearer ", "")
	signedCustomerID, err := h.tokenDecoder.Decode(token)
	if err != nil {
		writeErrorResponse(w, http.StatusForbidden, "Forbidden")
		return
	}

	inputs := mux.Vars(r)
	customerID, err := strconv.Atoi(inputs["customer_id"])
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if int64(customerID) != signedCustomerID {
		writeErrorResponse(w, http.StatusForbidden, "Forbidden")
		return
	}

	limit, offset := getLimitOffset(r.URL.Query())
	loans, err := h.loanRepo.View(ctx, int64(customerID), limit, offset)
	if err != nil && errors.Is(err, constants.ErrorRecordNotFound) {
		writeErrorResponse(w, http.StatusNotFound, "Not Found")
		return
	}
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	writeOKResponse(w, loans)
}
