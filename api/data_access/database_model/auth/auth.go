package auth

import "time"

//Auth represents the table definition of the Auths table in the db
type Auth struct {
	ID string
	//Sensitive information such as card details should be stored in compliance with PCI DSS requirement
	Number           string
	ExpiryDate       string
	AuthorisedAmount float32
	AvailableAmount  float32
	Currency         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}
