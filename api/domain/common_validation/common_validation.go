package common_validation

import (
	"log"
	"payment-gateway-api/api/config"
	"regexp"
	"time"
)

//IsValidUUID checks whether the field is in the UUID format
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile(config.UUIDCodeLayout)
	return r.MatchString(uuid)
}

//IsExpiryDateValid checks that the card is not expired
func IsExpiryDateValid(expiryDate string) bool {
	expirationDate, err := time.Parse(config.ExpirationDateLayout, expiryDate)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	currentTime, err := time.Parse(config.ExpirationDateLayout, time.Now().Format(config.ExpirationDateLayout))
	if err != nil {
		log.Println(err.Error())
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
