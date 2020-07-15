package void_service

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"payment-gateway-api/api/const/error_constant"
	"payment-gateway-api/api/data_access"
	"payment-gateway-api/api/data_access/database_model/auth"
	"payment-gateway-api/api/data_access/database_model/operation"
	"payment-gateway-api/api/domain/void_domain"
	"payment-gateway-api/api/services/common_service"
	"testing"
)

var (
	getAuthRecordByID        func(string) (bool, *auth.Auth, error)
	softDeleteAuthRecordByID func(string) error
	isAuthorisedState        func(string, string) (bool, error)
)

type databaseMock struct{}
type commonServiceMock struct{}

func (c commonServiceMock) IsAuthorisedState(operationName string, id string) (bool, error) {
	return isAuthorisedState(operationName, id)
}

func (d databaseMock) CheckRejectByCardNumber(string, string) (bool, error) {
	return true, nil
}

func (d databaseMock) UpdateAvailableAmountByAuthID(string, float32, string) error {
	return nil
}

func (d databaseMock) Setup(string) error {
	return nil
}

func (d databaseMock) InsertAuthRecord(*auth.Auth) error {
	return nil
}

func (d databaseMock) GetAuthRecordByID(id string) (bool, *auth.Auth, error) {
	return getAuthRecordByID(id)
}

func (d databaseMock) Close() error {
	return nil
}

func (d databaseMock) SoftDeleteAuthRecordByID(id string) error {
	return softDeleteAuthRecordByID(id)
}

func (d databaseMock) HardDeleteAuthRecordByID(string) error {
	return nil
}

func (d databaseMock) GetOperationByAuthIDAndOperationName(string, string) (bool, operation.Operation, error) {
	return false, operation.Operation{}, nil
}

func (d databaseMock) DeleteOperationRecordsByAuthID(string) error {
	return nil
}

func TestVoidService_VoidTransaction_NotVoidable(t *testing.T) {

	request := void_domain.VoidRequest{AuthId: "fc958d27-8e8e-4825-b3ec-e5236a8e7d28"}

	err1 := errors.New(error_constant.TransactionStateInvalid)
	expectedErrors := make([]error, 0)
	expectedErrors = append(expectedErrors, err1)

	isAuthorisedState = func(operationName string, id string) (bool, error) {
		return false, nil
	}

	common_service.CommonService = &commonServiceMock{}

	actualResponse, err := VoidService.VoidTransaction(request)
	assert.Nil(t, actualResponse)
	assert.EqualValues(t, fmt.Sprintf("%v", expectedErrors), err.ErrorMessage())
}

func TestVoidService_VoidTransaction(t *testing.T) {
	request := void_domain.VoidRequest{AuthId: "fc958d27-8e8e-4825-b3ec-e5236a8e7d28"}

	expectedResponse := void_domain.VoidResponse{
		IsSuccess: true,
		Amount:    10,
		Currency:  "GBP",
	}

	getAuthRecordByID = func(id string) (bool, *auth.Auth, error) {
		return true, &auth.Auth{
			AuthorisedAmount: expectedResponse.Amount,
			Currency:         expectedResponse.Currency,
		}, nil
	}
	isAuthorisedState = func(opName, id string) (b bool, err error) {
		return true, nil
	}

	softDeleteAuthRecordByID = func(s string) error {
		return nil
	}

	common_service.CommonService = &commonServiceMock{}
	data_access.Db = &databaseMock{}

	actualResponse, err := VoidService.VoidTransaction(request)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedResponse.IsSuccess, actualResponse.IsSuccess)
	assert.EqualValues(t, expectedResponse.Amount, actualResponse.Amount)
	assert.EqualValues(t, expectedResponse.Currency, actualResponse.Currency)
}

func TestVoidService_VoidTransactionAmount_GetAuthRecordError(t *testing.T) {
	request := void_domain.VoidRequest{
		AuthId: "fc958d27-8e8e-4825-b3ec-e5236a8e7d28",
	}

	err1 := errors.New(error_constant.TransactionNotFound)
	expectedErrors := make([]error, 0)
	expectedErrors = append(expectedErrors, err1)

	getAuthRecordByID = func(id string) (bool, *auth.Auth, error) {
		return true, &auth.Auth{}, errors.New("record not found")
	}

	isAuthorisedState = func(opName, id string) (b bool, err error) {
		return true, nil
	}

	common_service.CommonService = &commonServiceMock{}
	data_access.Db = &databaseMock{}

	actualResponse, err := VoidService.VoidTransaction(request)
	assert.Nil(t, actualResponse)
	assert.EqualValues(t, fmt.Sprintf("%v", expectedErrors), err.ErrorMessage())
}
