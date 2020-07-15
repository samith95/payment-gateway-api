package common_service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"payment-gateway-api/api/const/error_constant"
	"payment-gateway-api/api/data_access"
	"payment-gateway-api/api/data_access/database_model/auth"
	"payment-gateway-api/api/data_access/database_model/operation"
	"testing"
)

var (
	getOperationByAuthIDAndOperationName func(string, string) (bool, operation.Operation, error)
)

type databaseMock struct{}

func (d databaseMock) Setup(string) error {
	return nil
}

func (d databaseMock) InsertAuthRecord(*auth.Auth) error {
	return nil
}

func (d databaseMock) GetAuthRecordByID(string) (bool, *auth.Auth, error) {
	return true, &auth.Auth{}, nil
}

func (d databaseMock) Close() error {
	return nil
}

func (d databaseMock) SoftDeleteAuthRecordByID(string) error {
	return nil
}

func (d databaseMock) HardDeleteAuthRecordByID(string) error {
	return nil
}

func (d databaseMock) DeleteOperationRecordsByAuthID(string) error {
	return nil
}

func (d databaseMock) CheckRejectByCardNumber(string, string) (bool, error) {
	return true, nil
}

func (d databaseMock) UpdateAvailableAmountByAuthID(string, float32, string) error {
	return nil
}

func (d databaseMock) GetOperationByAuthIDAndOperationName(id string, opName string) (bool, operation.Operation, error) {
	return getOperationByAuthIDAndOperationName(id, opName)
}

func TestCommonService_IsAuthorisedState(t *testing.T) {
	getOperationByAuthIDAndOperationName = func(s string, s2 string) (b bool, o operation.Operation, err error) {
		return false, operation.Operation{}, nil
	}

	data_access.Db = &databaseMock{}

	isValid, err := CommonService.IsAuthorisedState("void", "valid_id")
	assert.Nil(t, err)
	assert.EqualValues(t, true, isValid)
}

func TestCommonService_IsAuthorisedState_NotAuthorised(t *testing.T) {
	getOperationByAuthIDAndOperationName = func(s string, s2 string) (b bool, o operation.Operation, err error) {
		return true, operation.Operation{}, nil
	}

	data_access.Db = &databaseMock{}

	isValid, err := CommonService.IsAuthorisedState("capture", "valid_id")
	assert.Nil(t, err)
	assert.EqualValues(t, false, isValid)
}

func TestCommonService_IsAuthorisedState_InvalidOperation(t *testing.T) {
	getOperationByAuthIDAndOperationName = func(s string, s2 string) (b bool, o operation.Operation, err error) {
		return true, operation.Operation{}, nil
	}

	data_access.Db = &databaseMock{}

	isValid, err := CommonService.IsAuthorisedState("invalid_operation", "valid_id")
	assert.EqualValues(t, error_constant.OperationNameInvalid, err.Error())
	assert.EqualValues(t, false, isValid)
}

func TestCommonService_IsAuthorisedState_ErrorFromDb(t *testing.T) {
	expectedError := "error"
	getOperationByAuthIDAndOperationName = func(s string, s2 string) (b bool, o operation.Operation, err error) {
		return false, operation.Operation{}, errors.New(expectedError)
	}

	data_access.Db = &databaseMock{}

	isValid, err := CommonService.IsAuthorisedState("void", "valid_id")
	assert.EqualValues(t, expectedError, err.Error())
	assert.EqualValues(t, false, isValid)
}
