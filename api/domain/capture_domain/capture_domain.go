package capture_domain

import (
	"errors"
	"payment-gateway-api/api/const/error_constant"
	"payment-gateway-api/api/domain/common_validation"
	"strings"
)

type CaptureRequest struct {
	AuthId string  `json:"id" binding:"required"`
	Amount float32 `json:"amount" binding:"required"`
}

type CaptureResponse struct {
	IsSuccess bool    `json:"success"`
	Amount    float32 `json:"amount"`
	Currency  string  `json:"currency"`
}

//ValidateFields strips all spaces from strings and checks their validity
func (v *CaptureRequest) ValidateFields() []error {
	var err = make([]error, 0)
	v.AuthId = strings.Replace(v.AuthId, " ", "", -1)
	if !common_validation.IsValidUUID(v.AuthId) {
		err = append(err, errors.New(error_constant.InvalidAuthIdField))
	}
	if !common_validation.IsAmountValid(v.Amount) {
		err = append(err, errors.New(error_constant.InvalidAmount))
	}
	return err
}
