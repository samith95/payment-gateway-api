package error_domain

import (
	"encoding/json"
	"fmt"
)

type GatewayErrorInterface interface {
	Status() int
	ErrorMessage() string
}

type GatewayError struct {
	Code  int
	Error string `json:"error"`
}

func (w *GatewayError) Status() int {
	return w.Code
}

func (e *GatewayError) ErrorMessage() string {
	return e.Error
}

func New(statusCode int, errorMsg ...error) GatewayErrorInterface {
	return &GatewayError{
		Code:  statusCode,
		Error: fmt.Sprintf("%v", errorMsg),
	}
}

func NewApiErrorFromBytes(body []byte) (GatewayErrorInterface, error) {
	var result GatewayError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
