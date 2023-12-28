package handlers

import (
	"aspire-lite/internals/models"
	"fmt"
	"net/http"
)

type RepaymentRepository interface {
	PayByCustomer(repayment *models.Repayment) error
}

type repayment struct {
	repaymentRepo RepaymentRepository
}

func NewRepayment(repaymentRepo RepaymentRepository) *repayment {
	return &repayment{repaymentRepo: repaymentRepo}
}

func (h *repayment) SubmitRepay(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}
