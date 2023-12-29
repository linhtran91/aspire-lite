package handlers

import (
	"aspire-lite/internals/models"
	"aspire-lite/internals/token"
	"aspire-lite/internals/usecases"
	mock "aspire-lite/mocks/handlers"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLogin(t *testing.T) {
	secretKey := "9876hjds"
	duration := 5 * time.Minute
	cases := []struct {
		name          string
		body          usecases.LoginInfo
		buildStubs    func(repo *mock.MockCustomerRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			body: usecases.LoginInfo{
				Username: "tester",
				Password: "321ndsjkqe",
			},
			buildStubs: func(repo *mock.MockCustomerRepository) {
				repo.EXPECT().GetUserCredential(gomock.Any(), gomock.Any()).Return(&models.Customer{
					ID:       1,
					Username: "tester",
					Password: "321ndsjkqe",
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "unauthorized",
			body: usecases.LoginInfo{
				Username: "tester",
				Password: "321ndsjkqe23",
			},
			buildStubs: func(repo *mock.MockCustomerRepository) {
				repo.EXPECT().GetUserCredential(gomock.Any(), gomock.Any()).Return(&models.Customer{
					ID:       1,
					Username: "tester",
					Password: "321ndsjkqe",
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "internal error",
			body: usecases.LoginInfo{
				Username: "tester",
				Password: "321ndsjkqe23",
			},
			buildStubs: func(repo *mock.MockCustomerRepository) {
				repo.EXPECT().GetUserCredential(gomock.Any(), gomock.Any()).Return(nil, errors.New("internal error"))
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

			repo := mock.NewMockCustomerRepository(ctrl)
			tokenBuilder := token.NewJWTTokenBuilder(secretKey, duration)
			tc.buildStubs(repo)

			handler := NewAuthenticator(repo, tokenBuilder)
			recorder := httptest.NewRecorder()
			data, _ := json.Marshal(tc.body)
			request, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(data))

			handler.Login(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
