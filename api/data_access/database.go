package data_access

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"payment-gateway-api/api/data_access/database_model/auth"
	"payment-gateway-api/api/data_access/database_model/operation"
	"time"
)

type database struct {
	Db *gorm.DB
}

type databaseInterface interface {
	Setup(string) error
	InsertAuthRecord(*auth.Auth) error
	GetAuthRecordByID(string) (bool, *auth.Auth, error)
	GetAllAuthRecords() ([]auth.Auth, error)
	Close() error
	SoftDeleteAuthRecordByID(string) error
	HardDeleteAuthRecordByID(string) error
	GetAllOperationsByAuthID(string) (bool, []operation.Operation, error)
	DeleteOperationRecordsByAuthID(string) error
}

var (
	Db databaseInterface = &database{}
)

func (db *database) Setup(dbStoreFilePath string) error {
	var err error
	//establish connection
	db.Db, err = gorm.Open("sqlite3", dbStoreFilePath)
	if err != nil {
		return err
	}

	//migrate struct definition into tables
	db.Db = db.Db.AutoMigrate(&auth.Auth{}, &operation.Operation{})
	if db.Db.Error != nil {
		return db.Db.Error
	}
	return nil
}

func (db *database) InsertAuthRecord(data *auth.Auth) error {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := insertOperation("authorisation", data, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func insertOperation(name string, data *auth.Auth, tx *gorm.DB) error {
	if err := tx.Error; err != nil {
		return err
	}

	op := &operation.Operation{
		AuthID:   data.ID,
		Name:     name,
		Amount:   data.AvailableAmount,
		Currency: data.Currency,
	}

	if err := tx.Create(op).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (db *database) GetAllAuthRecords() ([]auth.Auth, error) {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var records []auth.Auth
	if err := tx.Find(&records).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return records, nil
}

func (db *database) Close() error {
	return db.Db.Close()
}

func (db *database) GetAuthRecordByID(id string) (bool, *auth.Auth, error) {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var record auth.Auth
	if err := tx.Where("id = ?", id).First(&record).Error; err != nil {
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

func (db *database) SoftDeleteAuthRecordByID(id string) error {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var record auth.Auth
	if err := tx.Where("id = ?", id).First(&record).Error; err != nil {
		tx.Rollback()
		return err
	}

	record.DeletedAt = time.Now()

	if err := tx.Save(&record).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (db *database) HardDeleteAuthRecordByID(id string) error {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var record auth.Auth
	if err := tx.Where("id = ?", id).First(&record).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&record).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (db *database) DeleteOperationRecordsByAuthID(id string) error {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var record operation.Operation
	if err := tx.Where("auth_id = ?", id).Delete(&record).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (db *database) GetAllOperationsByAuthID(id string) (bool, []operation.Operation, error) {
	tx := db.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var records []operation.Operation

	if err := tx.Table("operations").Where("auth_id = ?", id).Find(&records).Error; err != nil {
		tx.Rollback()
		return false, nil, err
	}

	tx.Commit()
	return true, records, nil
}
