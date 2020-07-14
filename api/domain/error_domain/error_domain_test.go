package error_domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	err1 := errors.New("error1")
	err2 := errors.New("error2")
	errs := []error{err1, err2}
	code := 400
	actualError := New(code, errs...)
	assert.EqualValues(t, code, actualError.Status())
	assert.EqualValues(t, fmt.Sprintf("%v", errs), actualError.ErrorMessage())
}

func TestExchangeError(t *testing.T) {
	expectedError := GatewayError{
		Code:  400,
		Error: "Bad Request Error",
	}
	bytes, err := json.Marshal(expectedError)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	actualError, err := NewApiErrorFromBytes(bytes)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedError.Code, actualError.Status())
	assert.EqualValues(t, expectedError.Error, actualError.ErrorMessage())
}
