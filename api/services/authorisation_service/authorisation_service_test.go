package authorisation_service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"payment-gateway-api/api/const/error_constant"
	"payment-gateway-api/api/data_access"
	"payment-gateway-api/api/data_access/database_model/auth"
	"payment-gateway-api/api/data_access/database_model/operation"
	"payment-gateway-api/api/domain/auth_domain"
	"payment-gateway-api/api/domain/error_domain"
	"testing"
)

var (
	insertAuthRecord        func(*auth.Auth) error
	checkRejectByCardNumber func(string, string) (bool, error)
)

type databaseMock struct{}

func (db *databaseMock) GetOperationByAuthIDAndOperationName(string, string) (bool, operation.Operation, error) {
	return true, operation.Operation{}, nil
}

func (db *databaseMock) CheckRejectByCardNumber(opName string, cardNumber string) (bool, error) {
	return checkRejectByCardNumber(opName, cardNumber)
}

func (db *databaseMock) UpdateAvailableAmountByAuthID(string, float32, string) error {
	return nil
}

func (db *databaseMock) SoftDeleteAuthRecordByID(string) error {
	return nil
}

func (db *databaseMock) HardDeleteAuthRecordByID(string) error {
	return nil
}

func (db *databaseMock) DeleteOperationRecordsByAuthID(string) error {
	return nil
}

func (db *databaseMock) Setup(string) error {
	return nil
}

func (db *databaseMock) GetAuthRecordByID(string) (bool, *auth.Auth, error) {
	return false, nil, nil
}

func (db *databaseMock) Close() error {
	return nil
}

func (db *databaseMock) InsertAuthRecord(data *auth.Auth) error {
	return insertAuthRecord(data)
}

func TestAuthorisationService_AuthorisePayment(t *testing.T) {
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

	insertAuthRecord = func(auth *auth.Auth) error {
		return nil
	}

	checkRejectByCardNumber = func(opName string, cardNumber string) (b bool, err error) {
		return false, nil
	}

	data_access.Db = &databaseMock{}

	actualResponse, err := AuthorisationService.AuthoriseTransaction(request)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedResponse.IsSuccess, actualResponse.IsSuccess)
	assert.EqualValues(t, expectedResponse.Amount, actualResponse.Amount)
	assert.EqualValues(t, expectedResponse.Currency, actualResponse.Currency)
}

func TestAuthorisationService_AuthorisePayment_Error(t *testing.T) {
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

	insertAuthRecord = func(auth *auth.Auth) error {
		return errors.New(errorMessage)
	}

	data_access.Db = &databaseMock{}

	resp, actualError := AuthorisationService.AuthoriseTransaction(request)
	assert.Nil(t, resp)
	assert.EqualValues(t, expectedError.Status(), actualError.Status())
	assert.EqualValues(t, expectedError.ErrorMessage(), actualError.ErrorMessage())
}

func TestAuthorisationService_AuthorisePayment_RejectedCardError(t *testing.T) {

	request := auth_domain.AuthRequest{
		CardDetails: auth_domain.CardDetails{
			Number:     "4929907390318794",
			ExpiryDate: "12-2999",
			Cvv:        "123",
		},
		Amount:   10,
		Currency: "LKR",
	}

	checkRejectByCardNumber = func(opName string, cardNumber string) (bool, error) {
		return false, errors.New("")
	}

	data_access.Db = &databaseMock{}

	actualResponse, err := AuthorisationService.AuthoriseTransaction(request)
	assert.Nil(t, actualResponse)
	assert.EqualValues(t, error_constant.RejectRetrievalFailure, err.ErrorMessage())
}
