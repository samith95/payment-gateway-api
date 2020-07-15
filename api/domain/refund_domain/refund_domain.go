package refund_domain

import (
	"errors"
	"payment-gateway-api/api/const/error_constant"
	"payment-gateway-api/api/domain/common_validation"
	"strings"
)

//RefundRequest is the format for the request by the refund endpoint
type RefundRequest struct {
	AuthId string  `json:"id" binding:"required"`
	Amount float32 `json:"amount" binding:"required"`
}

//RefundResponse is the format for the response by the refund endpoint
type RefundResponse struct {
	IsSuccess bool    `json:"success"`
	Amount    float32 `json:"amount"`
	Currency  string  `json:"currency"`
}

//ValidateFields strips all spaces from strings and checks their validity
func (r *RefundRequest) ValidateFields() []error {
	var err = make([]error, 0)
	r.AuthId = strings.Replace(r.AuthId, " ", "", -1)
	if !common_validation.IsValidUUID(r.AuthId) {
		err = append(err, errors.New(error_constant.InvalidAuthIdField))
	}
	if !common_validation.IsAmountValid(r.Amount) {
		err = append(err, errors.New(error_constant.InvalidAmount))
	}
	return err
}
