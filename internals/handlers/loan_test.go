package handlers

import (
	"aspire-lite/internals/constants"
	"aspire-lite/internals/models"
	"aspire-lite/internals/token"
	"aspire-lite/internals/usecases"
	mock "aspire-lite/mocks/handlers"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateLoan(t *testing.T) {
	cases := []struct {
		name          string
		customerID    string
		body          usecases.Loan
		buildStubs    func(loanRepo *mock.MockLoanRepository, customerRepo *mock.MockCustomerRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "success",
			customerID: "1",
			body: usecases.Loan{
				Amount: 10000,
				Term:   3,
				Date:   "2022-02-08",
			},
			buildStubs: func(loanRepo *mock.MockLoanRepository, customerRepo *mock.MockCustomerRepository) {
				loanRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "bad request due to wrong date",
			customerID: "1",
			body: usecases.Loan{
				Amount: 10000,
				Term:   3,
				Date:   "2022-22-08",
			},
			buildStubs: func(loanRepo *mock.MockLoanRepository, customerRepo *mock.MockCustomerRepository) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "internal error",
			customerID: "1",
			body: usecases.Loan{
				Amount: 10000,
				Term:   3,
				Date:   "2022-02-08",
			},
			buildStubs: func(loanRepo *mock.MockLoanRepository, customerRepo *mock.MockCustomerRepository) {
				loanRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(-1), errors.New("internal error"))
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

			repo := mock.NewMockLoanRepository(ctrl)
			decoder := mock.NewMockTokenDecoder(ctrl)
			customerRepo := mock.NewMockCustomerRepository(ctrl)
			tc.buildStubs(repo, customerRepo)

			handler := NewLoan(repo, decoder, customerRepo)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/customers/%s/loans", tc.customerID)

			data, _ := json.Marshal(tc.body)
			request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			request = mux.SetURLVars(request, map[string]string{
				"customer_id": tc.customerID,
			})

			handler.CreateLoan(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestApproveLoan(t *testing.T) {
	cases := []struct {
		name          string
		loanID        string
		buildStubs    func(loanRepo *mock.MockLoanRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "success",
			loanID: "1",
			buildStubs: func(loanRepo *mock.MockLoanRepository) {
				loanRepo.EXPECT().Approve(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "internal error",
			loanID: "1",
			buildStubs: func(loanRepo *mock.MockLoanRepository) {
				loanRepo.EXPECT().Approve(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("internal error"))
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

			repo := mock.NewMockLoanRepository(ctrl)
			decoder := mock.NewMockTokenDecoder(ctrl)
			customerRepo := mock.NewMockCustomerRepository(ctrl)
			tc.buildStubs(repo)

			handler := NewLoan(repo, decoder, customerRepo)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/loans/%s/approve", tc.loanID)

			request, _ := http.NewRequest(http.MethodPut, url, nil)
			request = mux.SetURLVars(request, map[string]string{
				"loan_id": tc.loanID,
			})

			handler.ApproveLoan(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListLoan(t *testing.T) {
	secretKey := "9876hjds"
	duration := 5 * time.Minute
	cases := []struct {
		name          string
		customerID    int64
		page          int
		size          int
		setAuth       func(req *http.Request, encoder TokenEncoder, id int64)
		buildStubs    func(loanRepo *mock.MockLoanRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "success",
			customerID: 1,
			page:       1,
			size:       10,
			setAuth: func(req *http.Request, encoder TokenEncoder, id int64) {
				addToken(req, id, encoder, constants.AuthorizationKey, constants.AuthorizationHeader)
			},
			buildStubs: func(loanRepo *mock.MockLoanRepository) {
				loanRepo.EXPECT().View(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*models.Loan{
					{
						ID:         2,
						Amount:     3,
						Term:       10000,
						CustomerID: 1,
						Status:     1,
					},
					{
						ID:         4,
						Amount:     3,
						Term:       10000,
						CustomerID: 1,
						Status:     2,
					},
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				assert.NoError(t, err)

				var b JsonResponse
				err = json.Unmarshal(data, &b)
				assert.NoError(t, err)
				// log.Fatal(b.Data)
			},
		},
		{
			name:       "unauthorized",
			customerID: 1,
			setAuth:    func(req *http.Request, encoder TokenEncoder, id int64) {},
			buildStubs: func(loanRepo *mock.MockLoanRepository) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name:       "not found",
			customerID: 1,
			setAuth: func(req *http.Request, encoder TokenEncoder, id int64) {
				addToken(req, id, encoder, constants.AuthorizationKey, constants.AuthorizationHeader)
			},
			buildStubs: func(loanRepo *mock.MockLoanRepository) {
				loanRepo.EXPECT().View(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, constants.ErrorRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "internal error",
			customerID: 1,
			setAuth: func(req *http.Request, encoder TokenEncoder, id int64) {
				addToken(req, id, encoder, constants.AuthorizationKey, constants.AuthorizationHeader)
			},
			buildStubs: func(loanRepo *mock.MockLoanRepository) {
				loanRepo.EXPECT().View(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("internal error"))
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

			repo := mock.NewMockLoanRepository(ctrl)
			tokenBuilder := token.NewJWTTokenBuilder(secretKey, duration)
			customerRepo := mock.NewMockCustomerRepository(ctrl)
			tc.buildStubs(repo)

			handler := NewLoan(repo, tokenBuilder, customerRepo)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/customers/%d/loans", tc.customerID)

			request, _ := http.NewRequest(http.MethodGet, url, nil)
			request = mux.SetURLVars(request, map[string]string{
				"customer_id": strconv.Itoa(int(tc.customerID)),
			})

			tc.setAuth(request, tokenBuilder, tc.customerID)
			handler.List(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func addToken(r *http.Request, id int64, encoder TokenEncoder, authorizationType, authorKey string) {
	token, _ := encoder.Encode(id)
	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	r.Header.Set(authorKey, authorizationHeader)
}
