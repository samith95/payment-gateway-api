package void_domain

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVoidResponse(t *testing.T) {
	expectedResponse := VoidResponse{
		IsSuccess: true,
		Amount:    10,
		Currency:  "LKR",
	}

	bytes, err := json.Marshal(expectedResponse)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var actualResponse VoidResponse

	err = json.Unmarshal(bytes, &actualResponse)
	assert.Nil(t, err)
	assert.NotNil(t, actualResponse)
	assert.EqualValues(t, expectedResponse, actualResponse)
}

func TestVoidRequest_ValidateFields_Invalid(t *testing.T) {
	request := VoidRequest{
		AuthId: "invalid_id",
	}

	expectedErrors := []error{}
	expectedErrors = append(expectedErrors, errors.New("authorisation id field is not valid"))

	actualErrors := request.ValidateFields()

	assert.EqualValues(t, expectedErrors, actualErrors)
}

func TestVoidRequest_ValidateFields_Valid(t *testing.T) {
	request := VoidRequest{
		AuthId: "970c8844-9238-4c31-95ca-6f079dd65729",
	}

	actualErrors := request.ValidateFields()

	assert.EqualValues(t, []error{}, actualErrors)
}
