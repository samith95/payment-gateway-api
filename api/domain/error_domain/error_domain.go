package error_domain

import (
	"encoding/json"
	"fmt"
)

//GatewayErrorInterface is the used to interact with service errors
type GatewayErrorInterface interface {
	Status() int
	ErrorMessage() string
}

//GatewayError is the format for the error responses from the gateway
type GatewayError struct {
	Code  int    `json:"-"`
	Error string `json:"error"`
}

//Status returns the error code
func (w *GatewayError) Status() int {
	return w.Code
}

//ErrorMessage returns the error message
func (e *GatewayError) ErrorMessage() string {
	return e.Error
}

//New creates an error response struct given a collection error messages
func New(statusCode int, errorMsg ...error) GatewayErrorInterface {
	return &GatewayError{
		Code:  statusCode,
		Error: fmt.Sprintf("%v", errorMsg),
	}
}

//NewApiErrorFromBytes creates an error response struct given a byte array
func NewApiErrorFromBytes(body []byte) (GatewayErrorInterface, error) {
	var result GatewayError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
