package data_access

import (
	"github.com/stretchr/testify/assert"
	"payment-gateway-api/api/data_access/database_model/auth"
	"testing"
	"time"
)

func InitTestDb(t *testing.T) {
	err := Db.Setup("./test_db_store/test_gateway.db")
	assert.Nil(t, err)
}

func cleanupDB(id string, t *testing.T) {
	err := Db.HardDeleteAuthRecordByID(id)
	assert.Nil(t, err)

	err = Db.DeleteOperationRecordsByAuthID(id)
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

	//check that there are no operations saved
	isPresent, _, err := Db.GetOperationByAuthIDAndOperationName(expectedRecord.ID, "authorisation")
	assert.Nil(t, err)
	assert.EqualValues(t, false, isPresent)

	err = Db.InsertAuthRecord(&expectedRecord)

	//check that the authorisation operation has been saved
	isPresent, _, err = Db.GetOperationByAuthIDAndOperationName(expectedRecord.ID, "authorisation")
	assert.Nil(t, err)
	assert.EqualValues(t, true, isPresent)

	assert.Nil(t, err)

	_, actualRecord, err := Db.GetAuthRecordByID(expectedRecord.ID)

	assert.Nil(t, err)
	assert.EqualValues(t, expectedRecord.ID, actualRecord.ID)
	assert.EqualValues(t, expectedRecord.Number, actualRecord.Number)
	assert.EqualValues(t, expectedRecord.Currency, actualRecord.Currency)
	assert.EqualValues(t, expectedRecord.ExpiryDate, actualRecord.ExpiryDate)
	assert.EqualValues(t, expectedRecord.AuthorisedAmount, actualRecord.AuthorisedAmount)
	assert.EqualValues(t, expectedRecord.AvailableAmount, actualRecord.AvailableAmount)

	cleanupDB(expectedRecord.ID, t)
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

	cleanupDB(expectedRecord.ID, t)
}

func TestDatabase_CheckRejectByCardNumber(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	InitTestDb(t)
	defer Db.Close()

	rejectedCardNumber := "4000000000000119"
	nonRejectedCardNumber := "123"

	isPresent, err := Db.CheckRejectByCardNumber("authorisation", nonRejectedCardNumber)
	assert.Nil(t, err)
	assert.EqualValues(t, false, isPresent)

	isPresent, err = Db.CheckRejectByCardNumber("authorisation", rejectedCardNumber)
	assert.Nil(t, err)
	assert.EqualValues(t, true, isPresent)
}

func TestDatabase_UpdateAvailableAmountByAuthID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	InitTestDb(t)
	defer Db.Close()

	existingRecord := auth.Auth{
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

	expectedRecord := &auth.Auth{
		ID:               "NewCode",
		Number:           "123456789123456",
		ExpiryDate:       "12-2999",
		AuthorisedAmount: 10,
		AvailableAmount:  5,
		Currency:         "LKR",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		DeletedAt:        time.Time{},
	}

	err := Db.InsertAuthRecord(expectedRecord)
	assert.Nil(t, err)

	err = Db.UpdateAvailableAmountByAuthID(existingRecord.ID, expectedRecord.AvailableAmount, "capture")
	assert.Nil(t, err)

	_, actualRecord, err := Db.GetAuthRecordByID(expectedRecord.ID)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedRecord.ID, actualRecord.ID)
	assert.EqualValues(t, expectedRecord.Number, actualRecord.Number)
	assert.EqualValues(t, expectedRecord.Currency, actualRecord.Currency)
	assert.EqualValues(t, expectedRecord.ExpiryDate, actualRecord.ExpiryDate)
	assert.EqualValues(t, expectedRecord.AuthorisedAmount, actualRecord.AuthorisedAmount)
	assert.EqualValues(t, expectedRecord.AvailableAmount, actualRecord.AvailableAmount)

	cleanupDB(existingRecord.ID, t)
}
