package refund_service

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"payment-gateway-api/api/const/error_constant"
	"payment-gateway-api/api/data_access"
	"payment-gateway-api/api/data_access/database_model/auth"
	"payment-gateway-api/api/data_access/database_model/operation"
	"payment-gateway-api/api/domain/refund_domain"
	"payment-gateway-api/api/services/common_service"
	"testing"
)

var (
	getAuthRecordByID             func(string) (bool, *auth.Auth, error)
	softDeleteAuthRecordByID      func(string) error
	isAuthorisedState             func(string, string) (bool, error)
	checkRejectByCardNumber       func(string, string) (bool, error)
	updateAvailableAmountByAuthID func(string, float32, string) error
)

type databaseMock struct{}
type commonServiceMock struct{}

func (c commonServiceMock) IsAuthorisedState(operationName string, id string) (bool, error) {
	return isAuthorisedState(operationName, id)
}

func (d databaseMock) CheckRejectByCardNumber(opName string, cardNumber string) (bool, error) {
	return checkRejectByCardNumber(opName, cardNumber)
}

func (d databaseMock) UpdateAvailableAmountByAuthID(id string, newAmount float32, opName string) error {
	return updateAvailableAmountByAuthID(id, newAmount, opName)
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

func TestRefundService_RefundTransactionAmount_InvalidState(t *testing.T) {

	request := refund_domain.RefundRequest{
		AuthId: "fc958d27-8e8e-4825-b3ec-e5236a8e7d28",
		Amount: 10,
	}

	err1 := errors.New(error_constant.TransactionStateInvalid)
	expectedErrors := make([]error, 0)
	expectedErrors = append(expectedErrors, err1)

	isAuthorisedState = func(operationName string, id string) (bool, error) {
		return false, nil
	}

	common_service.CommonService = &commonServiceMock{}

	actualResponse, err := RefundService.RefundTransactionAmount(request)
	assert.Nil(t, actualResponse)
	assert.EqualValues(t, fmt.Sprintf("%v", expectedErrors), err.ErrorMessage())
}

func TestRefundService_RefundTransactionAmount(t *testing.T) {
	request := refund_domain.RefundRequest{
		AuthId: "fc958d27-8e8e-4825-b3ec-e5236a8e7d28",
		Amount: 5,
	}

	expectedResponse := refund_domain.RefundResponse{
		IsSuccess: true,
		Amount:    10,
		Currency:  "GBP",
	}

	capturedAmount := float32(5)

	getAuthRecordByID = func(id string) (bool, *auth.Auth, error) {
		return true, &auth.Auth{
			ExpiryDate:       "12-3999",
			AvailableAmount:  capturedAmount,
			AuthorisedAmount: capturedAmount + request.Amount,
			Currency:         expectedResponse.Currency,
		}, nil
	}

	checkRejectByCardNumber = func(opName string, cardNumber string) (bool, error) {
		return false, nil
	}

	updateAvailableAmountByAuthID = func(id string, newAmount float32, opName string) error {
		return nil
	}

	isAuthorisedState = func(opName, id string) (b bool, err error) {
		return true, nil
	}

	common_service.CommonService = &commonServiceMock{}
	data_access.Db = &databaseMock{}

	actualResponse, err := RefundService.RefundTransactionAmount(request)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedResponse.IsSuccess, actualResponse.IsSuccess)
	assert.EqualValues(t, expectedResponse.Amount, actualResponse.Amount)
	assert.EqualValues(t, expectedResponse.Currency, actualResponse.Currency)
}

func TestRefundService_RefundTransactionAmount_RejectedCardError(t *testing.T) {

	request := refund_domain.RefundRequest{
		AuthId: "fc958d27-8e8e-4825-b3ec-e5236a8e7d28",
		Amount: 5,
	}

	getAuthRecordByID = func(id string) (bool, *auth.Auth, error) {
		return true, &auth.Auth{}, nil
	}

	isAuthorisedState = func(opName, id string) (b bool, err error) {
		return true, nil
	}

	checkRejectByCardNumber = func(opName string, cardNumber string) (bool, error) {
		return false, errors.New("")
	}

	common_service.CommonService = &commonServiceMock{}
	data_access.Db = &databaseMock{}

	actualResponse, err := RefundService.RefundTransactionAmount(request)
	assert.Nil(t, actualResponse)
	assert.EqualValues(t, error_constant.RejectRetrievalFailure, err.ErrorMessage())
}

func TestRefundService_RefundTransactionAmount_UpdateAvailableAmountError(t *testing.T) {
	request := refund_domain.RefundRequest{
		AuthId: "fc958d27-8e8e-4825-b3ec-e5236a8e7d28",
		Amount: 5,
	}

	expectedResponse := refund_domain.RefundResponse{
		IsSuccess: true,
		Amount:    10,
		Currency:  "GBP",
	}

	capturedAmount := float32(5)

	getAuthRecordByID = func(id string) (bool, *auth.Auth, error) {
		return true, &auth.Auth{
			ExpiryDate:       "12-3999",
			AvailableAmount:  capturedAmount,
			AuthorisedAmount: capturedAmount + request.Amount,
			Currency:         expectedResponse.Currency,
		}, nil
	}

	checkRejectByCardNumber = func(opName string, cardNumber string) (bool, error) {
		return false, nil
	}

	updateAvailableAmountByAuthID = func(id string, newAmount float32, opName string) error {
		return errors.New("")
	}

	isAuthorisedState = func(opName, id string) (b bool, err error) {
		return true, nil
	}

	common_service.CommonService = &commonServiceMock{}
	data_access.Db = &databaseMock{}

	actualResponse, err := RefundService.RefundTransactionAmount(request)
	assert.Nil(t, actualResponse)
	assert.EqualValues(t, error_constant.UpdateAvailableAmountFailure, err.ErrorMessage())
}

func TestRefundService_RefundTransactionAmount_GetAuthRecordError(t *testing.T) {

	request := refund_domain.RefundRequest{
		AuthId: "fc958d27-8e8e-4825-b3ec-e5236a8e7d28",
		Amount: 5,
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

	actualResponse, err := RefundService.RefundTransactionAmount(request)
	assert.Nil(t, actualResponse)
	assert.EqualValues(t, fmt.Sprintf("%v", expectedErrors), err.ErrorMessage())
}
