package void_domain

import (
	"errors"
	"regexp"
	"strings"
)

type VoidRequest struct {
	AuthId string `json:"id" binding:"required"`
}

type VoidResponse struct {
	IsSuccess bool    `json:"success"`
	Amount    float32 `json:"amount"`
	Currency  string  `json:"currency"`
}

//ValidateFields strips all spaces from strings and checks their validity
func (v *VoidRequest) ValidateFields() []error {
	var err = make([]error, 0)
	v.AuthId = strings.Replace(v.AuthId, " ", "", -1)
	if !isValidUUID(v.AuthId) {
		err = append(err, errors.New("authorisation id field is not valid"))
	}
	return err
}

//isValidUUID checks whether the field is in the UUID format
func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
