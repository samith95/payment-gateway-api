package capture_domain

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"payment-gateway-api/api/const/error_constant"
	"testing"
)

func TestCaptureResponse(t *testing.T) {
	expectedResponse := CaptureResponse{
		IsSuccess: true,
		Amount:    10,
		Currency:  "LKR",
	}

	bytes, err := json.Marshal(expectedResponse)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var actualResponse CaptureResponse

	err = json.Unmarshal(bytes, &actualResponse)
	assert.Nil(t, err)
	assert.NotNil(t, actualResponse)
	assert.EqualValues(t, expectedResponse, actualResponse)
}

func TestCaptureRequest_ValidateFields_Invalid(t *testing.T) {
	request := CaptureRequest{
		AuthId: "invalid_id",
		Amount: 0,
	}

	expectedErrors := []error{}
	expectedErrors = append(expectedErrors, errors.New(error_constant.InvalidAuthIdField))
	expectedErrors = append(expectedErrors, errors.New(error_constant.InvalidAmount))

	actualErrors := request.ValidateFields()

	assert.EqualValues(t, expectedErrors, actualErrors)
}

func TestCaptureRequest_ValidateFields_Valid(t *testing.T) {
	request := CaptureRequest{
		AuthId: "970c8844-9238-4c31-95ca-6f079dd65729",
		Amount: 10,
	}

	actualErrors := request.ValidateFields()

	assert.EqualValues(t, []error{}, actualErrors)
}
