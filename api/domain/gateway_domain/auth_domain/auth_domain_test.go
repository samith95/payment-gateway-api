package auth_domain

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecommendationResponse(t *testing.T) {
	expectedResponse := AuthResponse{
		AuthID:    "123987-644ef1sdf-wf6d1fs1fr4w6f-df6ws54ef1",
		IsSuccess: true,
		Amount:    50,
		Currency:  "LKR",
	}

	bytes, err := json.Marshal(expectedResponse)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var actualResponse AuthResponse

	err = json.Unmarshal(bytes, &actualResponse)
	assert.Nil(t, err)
	assert.NotNil(t, actualResponse)
	assert.EqualValues(t, expectedResponse, actualResponse)
}

func TestAuthRequest_ValidateFields_Invalid(t *testing.T) {
	cardDetails := CardDetails{
		Number:     "4929907390318797",
		ExpiryDate: "01-1900",
		Cvv:        "78969",
	}
	request := AuthRequest{
		CardDetails: cardDetails,
		Amount:      -5,
		Currency:    "invalid-currency",
	}

	expectedErrors := []error{}
	expectedErrors = append(expectedErrors, errors.New("card number is not valid"))
	expectedErrors = append(expectedErrors, errors.New("expiry date is not valid"))
	expectedErrors = append(expectedErrors, errors.New("cvv number is not valid"))
	expectedErrors = append(expectedErrors, errors.New("amount cannot be negative"))
	expectedErrors = append(expectedErrors, errors.New("currency code is invalid"))

	errs := request.ValidateFields()

	assert.EqualValues(t, expectedErrors, errs)
}

func TestAuthRequest_ValidateFields_Valid(t *testing.T) {
	cardDetails := CardDetails{
		Number:     "4929907390318794",
		ExpiryDate: "12-3500",
		Cvv:        "123",
	}
	request := AuthRequest{
		CardDetails: cardDetails,
		Amount:      10000,
		Currency:    "GBP",
	}

	expectedErrors := []error{}

	errs := request.ValidateFields()

	assert.EqualValues(t, expectedErrors, errs)
}
