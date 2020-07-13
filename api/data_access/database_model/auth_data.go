package database_model

import (
	"time"
)

//AuditData will be used to manage audit information about the various table definitions.
type AuditData struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type Auth struct {
	ID string
	//Sensitive information such as card details should be stored in compliance with PCI DSS requirement
	Number           string
	ExpiryDate       string
	AuthorisedAmount float32
	AvailableAmount  float32
	Currency         string
	AuditData        AuditData
}
