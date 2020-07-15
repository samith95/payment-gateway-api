package auth_domain

import (
	"errors"
	"github.com/joeljunstrom/go-luhn"
	"payment-gateway-api/api/const/error_constant"
	"payment-gateway-api/api/domain/common_validation"
	"regexp"
	"strings"
)

type AuthRequest struct {
	CardDetails CardDetails `json:"card_details" binding:"required"`
	Amount      float32     `json:"amount" binding:"required"`
	Currency    string      `json:"currency" binding:"required"`
}

type CardDetails struct {
	Number     string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	Cvv        string `json:"cvv"`
}

type AuthResponse struct {
	AuthID    string  `json:"id"`
	IsSuccess bool    `json:"success"`
	Amount    float32 `json:"amount"`
	Currency  string  `json:"currency"`
}

//ValidateFields strips all spaces from strings and checks their validity
func (r *AuthRequest) ValidateFields() []error {
	var err = make([]error, 0)
	r.CardDetails.Number = strings.Replace(r.CardDetails.Number, " ", "", -1)
	if !isCardNumberValid(r.CardDetails.Number) {
		err = append(err, errors.New(error_constant.InvalidCardNumber))
	}
	r.CardDetails.ExpiryDate = strings.Replace(r.CardDetails.ExpiryDate, " ", "", -1)
	if !common_validation.IsExpiryDateValid(r.CardDetails.ExpiryDate) {
		err = append(err, errors.New(error_constant.InvalidCardExpiryDate))
	}
	r.CardDetails.Cvv = strings.Replace(r.CardDetails.Cvv, " ", "", -1)
	if !isCvvValid(r.CardDetails.Cvv) {
		err = append(err, errors.New(error_constant.InvalidCvv))
	}
	if !common_validation.IsAmountValid(r.Amount) {
		err = append(err, errors.New(error_constant.InvalidAmount))
	}
	if !isCurrencyCodeValid(r.Currency) {
		err = append(err, errors.New(error_constant.InvalidCurrencyCode))
	}
	return err
}

//isCardNumberValid checks the card number validity using the Luhn algorithm
func isCardNumberValid(cardNumber string) bool {
	return luhn.Valid(cardNumber)
}

//isCvvValid checks that the CVV is made of 3 or 4 integers
func isCvvValid(cvv string) bool {
	isValid, _ := regexp.MatchString("^[0-9]{3,4}$", cvv)
	return isValid
}

//isCurrencyCodeValid checks the currency is a 3 letter string
func isCurrencyCodeValid(currency string) bool {
	isValid, _ := regexp.MatchString("^[A-Z]{3}$", currency)
	return isValid
}
