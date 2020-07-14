package data_access

import (
	"github.com/stretchr/testify/assert"
	"payment-gateway-api/api/data_access/database_model/auth"
	"payment-gateway-api/api/data_access/database_model/operation"
	"testing"
	"time"
)

func InitTestDb(t *testing.T) {
	err := Db.Setup("./test_db_store/test_gateway.db")
	assert.Nil(t, err)
}

func TestDatabase_CreateAuthRecord_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	InitTestDb(t)
	defer Db.Close()

	expectedRecord := auth.Auth{
		ID:               "NewCode",
		Number:           "123456789123456",
		ExpiryDate:       "12-2999",
		AuthorisedAmount: 10,
		AvailableAmount:  10,
		Currency:         "LKR",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		DeletedAt:        time.Time{},
	}

	ok, operations, err := Db.GetAllOperationsByAuthID(expectedRecord.ID)
	assert.Nil(t, err)
	assert.EqualValues(t, []operation.Operation{}, operations)
	assert.EqualValues(t, true, ok)

	err = Db.InsertAuthRecord(&expectedRecord)

	//check that the authorisation operation has been saved
	ok, operations, err = Db.GetAllOperationsByAuthID(expectedRecord.ID)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(operations))
	assert.EqualValues(t, true, ok)

	assert.Nil(t, err)

	_, actualRecord, err := Db.GetAuthRecordByID(expectedRecord.ID)

	assert.Nil(t, err)
	assert.EqualValues(t, expectedRecord.ID, actualRecord.ID)
	assert.EqualValues(t, expectedRecord.Number, actualRecord.Number)
	assert.EqualValues(t, expectedRecord.Currency, actualRecord.Currency)
	assert.EqualValues(t, expectedRecord.ExpiryDate, actualRecord.ExpiryDate)
	assert.EqualValues(t, expectedRecord.AuthorisedAmount, actualRecord.AuthorisedAmount)
	assert.EqualValues(t, expectedRecord.AvailableAmount, actualRecord.AvailableAmount)

	err = Db.HardDeleteAuthRecordByID(expectedRecord.ID)
	assert.Nil(t, err)

	err = Db.DeleteOperationRecordsByAuthID(expectedRecord.ID)
	assert.Nil(t, err)
}

func TestDatabase_SoftDeleteAuthRecordByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	InitTestDb(t)
	defer Db.Close()

	expectedRecord := auth.Auth{
		ID:               "NewCode",
		Number:           "123456789123456",
		ExpiryDate:       "12-2999",
		AuthorisedAmount: 10,
		AvailableAmount:  10,
		Currency:         "LKR",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		DeletedAt:        time.Time{},
	}

	err := Db.InsertAuthRecord(&expectedRecord)
	assert.Nil(t, err)

	err = Db.SoftDeleteAuthRecordByID(expectedRecord.ID)
	assert.Nil(t, err)

	_, actualRecord, err := Db.GetAuthRecordByID(expectedRecord.ID)
	assert.Nil(t, err)
	assert.EqualValues(t, &auth.Auth{}, actualRecord)

	//cleanup
	err = Db.HardDeleteAuthRecordByID(expectedRecord.ID)
	assert.Nil(t, err)

	err = Db.DeleteOperationRecordsByAuthID(expectedRecord.ID)
	assert.Nil(t, err)
}
