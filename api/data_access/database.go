package data_access

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"payment-gateway-api/api/config"
	"payment-gateway-api/api/data_access/database_model"
)

type database struct {
	Db *gorm.DB
}

type databaseInterface interface {
	Init() error
	CreateAuthRecord(data *database_model.Auth) error
	GetAllRecords() ([]database_model.Auth, error)
	Close() error
}

var (
	Db databaseInterface = &database{}
)

func (db *database) Init() error {
	var err error
	//establish connection
	db.Db, err = gorm.Open("sqlite3", config.DbStoreFilePath)
	if err != nil {
		return err
	}

	//migrate struct definition into tables
	db.Db = db.Db.AutoMigrate(&database_model.Auth{})
	if db.Db.Error != nil {
		return db.Db.Error
	}
	return nil
}

func (db *database) CreateAuthRecord(data *database_model.Auth) error {
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

	return tx.Commit().Error
}

func (db *database) GetAllRecords() ([]database_model.Auth, error) {
	var records []database_model.Auth
	if err := db.Db.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (db *database) Close() error {
	return db.Db.Close()
}
