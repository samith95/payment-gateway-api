package data_access

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"payment-gateway-api/api/data_access/database_model/auth"
	"payment-gateway-api/api/data_access/database_model/operation"
	"payment-gateway-api/api/data_access/database_model/reject"
	"strings"
	"time"
)

type database struct {
	Db *gorm.DB
}

type databaseInterface interface {
	Setup(string) error
	InsertAuthRecord(*auth.Auth) error
	GetAuthRecordByID(string) (bool, *auth.Auth, error)
	Close() error
	SoftDeleteAuthRecordByID(string) error
	HardDeleteAuthRecordByID(string) error
	GetOperationByAuthIDAndOperationName(string, string) (bool, operation.Operation, error)
	DeleteOperationRecordsByAuthID(string) error
	CheckRejectByCardNumber(string, string) (bool, error)
	UpdateAvailableAmountByAuthID(string, float32, string) error
}

var (
	Db databaseInterface = &database{}
)

//Setup opens the db and the relevant tables
func (db *database) Setup(dbStoreFilePath string) error {
	var err error
	//establish connection
	db.Db, err = gorm.Open("sqlite3", dbStoreFilePath)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	//migrate struct definition into tables
	db.Db = db.Db.AutoMigrate(&auth.Auth{}, &operation.Operation{}, &reject.Reject{})
	if db.Db.Error != nil {
		return db.Db.Error
	}
	return nil
}

//InsertAuthRecord inserts an entry into the auths table
func (db *database) InsertAuthRecord(data *auth.Auth) error {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Println(err.Error())
		return err
	}

	if err := tx.Create(data).Error; err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}

	if err := insertOperation("authorisation", data, tx); err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func insertOperation(name string, data *auth.Auth, tx *gorm.DB) error {
	if err := tx.Error; err != nil {
		log.Println(err.Error())
		return err
	}

	op := &operation.Operation{
		AuthID:   data.ID,
		Name:     name,
		Amount:   data.AvailableAmount,
		Currency: data.Currency,
	}

	if err := tx.Create(op).Error; err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

//Close closes the connection to the db
func (db *database) Close() error {
	return db.Db.Close()
}

//GetAuthRecordByID fetches an auth record given its id
func (db *database) GetAuthRecordByID(id string) (bool, *auth.Auth, error) {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var record auth.Auth
	if err := tx.Where("id = ?", id).First(&record).Error; err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return false, nil, err
	}

	tx.Commit()

	//if the auth record has been soft deleted, return empty struct
	var empty time.Time
	if record.DeletedAt != empty {
		return false, &auth.Auth{}, nil
	}

	return true, &record, nil
}

//SoftDeleteAuthRecordByID initialises the deleteAt auth's variable
func (db *database) SoftDeleteAuthRecordByID(id string) error {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var record auth.Auth
	if err := tx.Where("id = ?", id).First(&record).Error; err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}

	record.DeletedAt = time.Now()

	if err := tx.Save(&record).Error; err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//HardDeleteAuthRecordByID removes the auth record given its id
func (db *database) HardDeleteAuthRecordByID(id string) error {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var record auth.Auth
	if err := tx.Where("id = ?", id).First(&record).Error; err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&record).Error; err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//DeleteOperationRecordsByAuthID removes all operations of a given authorisation id
func (db *database) DeleteOperationRecordsByAuthID(id string) error {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var record operation.Operation
	if err := tx.Where("auth_id = ?", id).Delete(&record).Error; err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

//GetOperationByAuthIDAndOperationName fetches the operation given the authorisation ID and the operation name to look for
func (db *database) GetOperationByAuthIDAndOperationName(id string, opName string) (bool, operation.Operation, error) {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var record operation.Operation

	err := tx.Table("operations").Where("auth_id = ? AND name = ?", id, opName).Find(&record).Error

	if err != nil && err.Error() != "record not found" {
		log.Println(err.Error())
		tx.Rollback()
		return false, operation.Operation{}, err
	}

	//if it is not found, then return empty struct
	emptyStruct := operation.Operation{}
	if record == emptyStruct {
		return false, emptyStruct, tx.Commit().Error
	}

	return true, record, tx.Commit().Error
}

//CheckRejectByCardNumber checks whether the operation with the passed card number is present in the rejects table
func (db *database) CheckRejectByCardNumber(operation string, cardNumber string) (bool, error) {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var record reject.Reject

	err := tx.Where("card_number = ?", cardNumber).First(&record).Error
	if err != nil && err.Error() != "record not found" {
		log.Println(err.Error())
		tx.Rollback()
		return false, err
	}

	//check if operation field contains current operation
	if strings.Contains(record.Operation, operation) {
		return true, tx.Commit().Error
	}

	return false, tx.Commit().Error
}

//UpdateAvailableAmountByAuthID updates the available amount of the given authorisation id record
func (db *database) UpdateAvailableAmountByAuthID(id string, amount float32, opName string) error {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var record auth.Auth
	if err := tx.Where("id = ?", id).First(&record).Error; err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}

	record.AvailableAmount = amount

	if err := tx.Model(&record).Where("id = ?", id).Update("available_amount", amount).Error; err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}

	if err := insertOperation(opName, &record, tx); err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
