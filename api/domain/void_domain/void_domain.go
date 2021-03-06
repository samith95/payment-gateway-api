package void_domain

import (
	"errors"
	"payment-gateway-api/api/const/error_constant"
	"payment-gateway-api/api/domain/common_validation"
	"strings"
)

//VoidRequest is the format for the request by the void endpoint
type VoidRequest struct {
	AuthId string `json:"id" binding:"required"`
}

//VoidResponse is the format for the response by the void endpoint
type VoidResponse struct {
	IsSuccess bool    `json:"success"`
	Amount    float32 `json:"amount"`
	Currency  string  `json:"currency"`
}

//ValidateFields strips all spaces from strings and checks their validity
func (v *VoidRequest) ValidateFields() []error {
	var err = make([]error, 0)
	v.AuthId = strings.Replace(v.AuthId, " ", "", -1)
	if !common_validation.IsValidUUID(v.AuthId) {
		err = append(err, errors.New(error_constant.InvalidAuthIdField))
	}
	return err
}
