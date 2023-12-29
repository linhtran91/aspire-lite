package handlers

import (
	"aspire-lite/internals/constants"
	"aspire-lite/internals/models"
	"aspire-lite/internals/usecases"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type LoanRepository interface {
	Create(ctx context.Context, loan *models.Loan, repayments []*models.Repayment) (int64, error)
	Approve(ctx context.Context, loanID int64, at time.Time) error
	View(ctx context.Context, customerID int64, limit, offset int) ([]*models.Loan, error)
	UpdateStatus(ctx context.Context, loanID int64) error
}

type loanHandler struct {
	loanRepo LoanRepository
}

func NewLoan(loanRepo LoanRepository) *loanHandler {
	return &loanHandler{loanRepo: loanRepo}
}

func (h *loanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var loan usecases.Loan
	if err := json.NewDecoder(r.Body).Decode(&loan); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}
	inputs := mux.Vars(r)
	customerID, err := strconv.Atoi(inputs["customer_id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Bad request")
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
	repayments := make([]*models.Repayment, len(loan.Repayments))
	for i, repay := range loan.Repayments {
		scheduledDate, err := parseDate(repay.Date)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "Bad request")
			return
		}
		repayments[i] = &models.Repayment{
			ID:              generateUUID(),
			ScheduledAmount: repay.Amount,
			Status:          constants.PENDING,
			ScheduledPayAt:  scheduledDate,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
	}

	id, err := h.loanRepo.Create(ctx, newLoan, repayments)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	writeOKResponse(w, id)
}

func (h *loanHandler) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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
	writeOKResponse(w, loanID)
}

func (h *loanHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	inputs := mux.Vars(r)
	customerID, err := strconv.Atoi(inputs["customer_id"])
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	limit, offset := getPageAndSize(r.URL.Query())
	loans, err := h.loanRepo.View(ctx, int64(customerID), limit, offset)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	writeOKResponse(w, loans)
}
