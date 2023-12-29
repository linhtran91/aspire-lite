package handlers

import (
	"aspire-lite/internals/constants"
	"aspire-lite/internals/models"
	"aspire-lite/internals/usecases"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type RepaymentRepository interface {
	SubmitRepayment(ctx context.Context, repayment *models.Repayment) error
	GetByID(ctx context.Context, id string) (*models.Repayment, error)
	CountUnpaidRepayment(ctx context.Context, loanID int64) (int64, error)
}

type repayment struct {
	repaymentRepo RepaymentRepository
	loanRepo      LoanRepository
}

func NewRepayment(repaymentRepo RepaymentRepository, loanRepo LoanRepository) *repayment {
	return &repayment{repaymentRepo: repaymentRepo, loanRepo: loanRepo}
}

func (h *repayment) SubmitRepay(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
	defer cancel()
	inputs := mux.Vars(r)
	id := inputs["repayment_id"]
	repayment, err := h.repaymentRepo.GetByID(ctx, id)
	if err != nil && errors.Is(err, constants.ErrorRecordNotFound) {
		writeErrorResponse(w, http.StatusNotFound, "Not Found")
		return
	}
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}

	var current usecases.SubmittedRepayment
	if err := json.NewDecoder(r.Body).Decode(&current); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}

	if repayment.ScheduledAmount > current.Amount {
		writeErrorResponse(w, http.StatusBadRequest, "Amount should be greater or equal to the scheduled repayment")
		return
	}

	now := time.Now().UTC()
	repayment.ActualAmount = current.Amount
	repayment.PaidAt = now
	repayment.Status = constants.PAID
	repayment.UpdatedAt = now
	if err := h.repaymentRepo.SubmitRepayment(ctx, repayment); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	count, err := h.repaymentRepo.CountUnpaidRepayment(ctx, repayment.LoanID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if count > 0 {
		writeOKResponse(w, map[string]interface{}{
			"repayment_id": id,
		})
		return
	}

	if err := h.loanRepo.UpdateStatus(ctx, repayment.LoanID, now); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	writeOKResponse(w, map[string]interface{}{
		"repayment_id": id,
		"loan_id":      repayment.LoanID,
	})
}
