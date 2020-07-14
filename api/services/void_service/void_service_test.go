package void_service

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"payment-gateway-api/api/data_access"
	"payment-gateway-api/api/data_access/database_model/auth"
	"payment-gateway-api/api/data_access/database_model/operation"
	"payment-gateway-api/api/domain/void_domain"
	"testing"
)

var (
	getAllOperationsByAuthID func(string) (bool, []operation.Operation, error)
	getAuthRecordByID        func(id string) (bool, *auth.Auth, error)
	softDeleteAuthRecordByID func(string) error
)

type databaseMock struct{}

func (d databaseMock) Setup(string) error {
	return nil
}

func (d databaseMock) InsertAuthRecord(*auth.Auth) error {
	panic("implement me")
}

func (d databaseMock) GetAuthRecordByID(id string) (bool, *auth.Auth, error) {
	return getAuthRecordByID(id)
}

func (d databaseMock) GetAllAuthRecords() ([]auth.Auth, error) {
	panic("implement me")
}

func (d databaseMock) Close() error {
	panic("implement me")
}

func (d databaseMock) SoftDeleteAuthRecordByID(id string) error {
	return softDeleteAuthRecordByID(id)
}

func (d databaseMock) HardDeleteAuthRecordByID(string) error {
	return nil
}

func (d databaseMock) GetAllOperationsByAuthID(id string) (bool, []operation.Operation, error) {
	return getAllOperationsByAuthID(id)
}

func (d databaseMock) DeleteOperationRecordsByAuthID(string) error {
	return nil
}

func TestVoidService_VoidTransaction_NotVoidable(t *testing.T) {

	request := void_domain.VoidRequest{AuthId: "fc958d27-8e8e-4825-b3ec-e5236a8e7d28"}

	err1 := errors.New("transaction is not in a state that allow cancellation")
	expectedErrors := make([]error, 0)
	expectedErrors = append(expectedErrors, err1)

	getAllOperationsByAuthID = func(s string) (bool, []operation.Operation, error) {
		return true, []operation.Operation{operation.Operation{}, operation.Operation{}}, nil
	}

	data_access.Db = &databaseMock{}

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

	getAllOperationsByAuthID = func(s string) (bool, []operation.Operation, error) {
		return true, []operation.Operation{operation.Operation{}}, nil
	}

	getAuthRecordByID = func(id string) (bool, *auth.Auth, error) {
		return true, &auth.Auth{
			AuthorisedAmount: expectedResponse.Amount,
			Currency:         expectedResponse.Currency,
		}, nil
	}

	softDeleteAuthRecordByID = func(s string) error {
		return nil
	}

	data_access.Db = &databaseMock{}

	actualResponse, err := VoidService.VoidTransaction(request)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedResponse.IsSuccess, actualResponse.IsSuccess)
	assert.EqualValues(t, expectedResponse.Amount, actualResponse.Amount)
	assert.EqualValues(t, expectedResponse.Currency, actualResponse.Currency)
}
