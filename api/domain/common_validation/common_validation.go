package common_validation

import (
	"payment-gateway-api/api/config"
	"regexp"
	"time"
)

//IsValidUUID checks whether the field is in the UUID format
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

//IsExpiryDateValid checks that the card is not expired
func IsExpiryDateValid(expiryDate string) bool {
	expirationDate, err := time.Parse(config.ExpirationDateLayout, expiryDate)
	if err != nil {
		return false
	}

	currentTime, err := time.Parse(config.ExpirationDateLayout, time.Now().Format(config.ExpirationDateLayout))
	if err != nil {
		return false
	}

	//check that card is not expired
	if currentTime.After(expirationDate) {
		return false
	}

	return true
}

//isAmountValid checks in case amount is negative or zero
func IsAmountValid(amount float32) bool {
	return amount > 0
}
