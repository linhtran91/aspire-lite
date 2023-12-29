package handlers

import (
	"aspire-lite/internals/constants"
	"aspire-lite/internals/models"
	"aspire-lite/internals/usecases"
	mock "aspire-lite/mocks/handlers"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSubmitRepay(t *testing.T) {
	cases := []struct {
		name          string
		repaymentID   string
		body          usecases.SubmittedRepayment
		buildStubs    func(loanRepo *mock.MockLoanRepository, repaymentRepo *mock.MockRepaymentRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "success without updating loan",
			repaymentID: "56yewew",
			body: usecases.SubmittedRepayment{
				Amount: 2500,
			},
			buildStubs: func(loanRepo *mock.MockLoanRepository, repaymentRepo *mock.MockRepaymentRepository) {
				repayment := &models.Repayment{
					ID:              "56yewew",
					LoanID:          1,
					ScheduledAmount: 2000,
					Status:          1,
				}
				repaymentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(repayment, nil)
				repaymentRepo.EXPECT().SubmitRepayment(gomock.Any(), gomock.Eq(repayment)).Return(nil)
				repaymentRepo.EXPECT().CountUnpaidRepayment(gomock.Any(), repayment.LoanID).Return(int64(2), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:        "success with updating loan",
			repaymentID: "56yewew",
			body: usecases.SubmittedRepayment{
				Amount: 2500,
			},
			buildStubs: func(loanRepo *mock.MockLoanRepository, repaymentRepo *mock.MockRepaymentRepository) {
				repayment := &models.Repayment{
					ID:              "56yewew",
					LoanID:          1,
					ScheduledAmount: 2000,
					Status:          1,
				}
				repaymentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(repayment, nil)
				repaymentRepo.EXPECT().SubmitRepayment(gomock.Any(), gomock.Eq(repayment)).Return(nil)
				repaymentRepo.EXPECT().CountUnpaidRepayment(gomock.Any(), repayment.LoanID).Return(int64(0), nil)
				loanRepo.EXPECT().UpdateStatus(gomock.Any(), int64(1), gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:        "fail due to lower summitted amount",
			repaymentID: "56yewew",
			body: usecases.SubmittedRepayment{
				Amount: 2500,
			},
			buildStubs: func(loanRepo *mock.MockLoanRepository, repaymentRepo *mock.MockRepaymentRepository) {
				repayment := &models.Repayment{
					ID:              "56yewew",
					LoanID:          1,
					ScheduledAmount: 3000,
					Status:          1,
				}
				repaymentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(repayment, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "not found repayment",
			repaymentID: "56yewew",
			body: usecases.SubmittedRepayment{
				Amount: 2500,
			},
			buildStubs: func(loanRepo *mock.MockLoanRepository, repaymentRepo *mock.MockRepaymentRepository) {
				repaymentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, constants.ErrorRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:        "fail due to internal error",
			repaymentID: "56yewew",
			body: usecases.SubmittedRepayment{
				Amount: 2500,
			},
			buildStubs: func(loanRepo *mock.MockLoanRepository, repaymentRepo *mock.MockRepaymentRepository) {
				repayment := &models.Repayment{
					ID:              "56yewew",
					LoanID:          1,
					ScheduledAmount: 2000,
					Status:          1,
				}
				repaymentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(repayment, nil)
				repaymentRepo.EXPECT().SubmitRepayment(gomock.Any(), gomock.Eq(repayment)).Return(nil)
				repaymentRepo.EXPECT().CountUnpaidRepayment(gomock.Any(), repayment.LoanID).Return(int64(-1), errors.New("internal error"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			loanRepo := mock.NewMockLoanRepository(ctrl)
			repaymentRepo := mock.NewMockRepaymentRepository(ctrl)
			tc.buildStubs(loanRepo, repaymentRepo)

			handler := NewRepayment(repaymentRepo, loanRepo)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/repayments/%s", tc.repaymentID)

			data, _ := json.Marshal(tc.body)
			request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			request = mux.SetURLVars(request, map[string]string{
				"repayment_id": tc.repaymentID,
			})

			handler.SubmitRepay(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
