package authorisation_service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"payment-gateway-api/api/data_access"
	"payment-gateway-api/api/data_access/database_model"
	"payment-gateway-api/api/domain/gateway_domain/auth_domain"
	"payment-gateway-api/api/domain/gateway_domain/error_domain"
	"testing"
)

var (
	createAuthRecord func(*database_model.Auth) error
)

type databaseMock struct{}

func (db *databaseMock) Init() error {
	return nil
}

func (db *databaseMock) GetAllRecords() ([]database_model.Auth, error) {
	return nil, nil
}

func (db *databaseMock) Close() error {
	return nil
}

func (db *databaseMock) CreateAuthRecord(data *database_model.Auth) error {
	return createAuthRecord(data)
}

func TestAuthorisationServiceAuthorisePayment(t *testing.T) {
	cardDetails := auth_domain.CardDetails{
		Number:     "4929907390318794",
		ExpiryDate: "12-3500",
		Cvv:        "123",
	}
	request := auth_domain.AuthRequest{
		CardDetails: cardDetails,
		Amount:      10000,
		Currency:    "GBP",
	}

	expectedResponse := auth_domain.AuthResponse{
		AuthID:    "",
		IsSuccess: true,
		Amount:    request.Amount,
		Currency:  request.Currency,
	}

	createAuthRecord = func(auth *database_model.Auth) error {
		return nil
	}

	data_access.Db = &databaseMock{}

	actualResponse, err := AuthorisationService.AuthorisePayment(request)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedResponse.IsSuccess, actualResponse.IsSuccess)
	assert.EqualValues(t, expectedResponse.Amount, actualResponse.Amount)
	assert.EqualValues(t, expectedResponse.Currency, actualResponse.Currency)
}

func TestAuthorisationServiceAuthorisePaymentError(t *testing.T) {
	cardDetails := auth_domain.CardDetails{
		Number:     "4929907390318794",
		ExpiryDate: "12-3500",
		Cvv:        "123",
	}
	request := auth_domain.AuthRequest{
		CardDetails: cardDetails,
		Amount:      10000,
		Currency:    "GBP",
	}

	errorMessage := "cannot connect to db"

	expectedError := error_domain.GatewayError{
		Code:  http.StatusInternalServerError,
		Error: errorMessage,
	}

	createAuthRecord = func(auth *database_model.Auth) error {
		return errors.New(errorMessage)
	}

	data_access.Db = &databaseMock{}

	resp, actualError := AuthorisationService.AuthorisePayment(request)
	assert.Nil(t, resp)
	assert.EqualValues(t, expectedError.Status(), actualError.Status())
	assert.EqualValues(t, expectedError.ErrorMessage(), actualError.ErrorMessage())
}
