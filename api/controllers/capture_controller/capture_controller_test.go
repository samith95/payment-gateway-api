package capture_controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"payment-gateway-api/api/domain/capture_domain"
	"payment-gateway-api/api/domain/error_domain"
	"payment-gateway-api/api/services/capture_service"
	"strings"
	"testing"
)

var (
	captureTransactionAmount func(request capture_domain.CaptureRequest) (*capture_domain.CaptureResponse, error_domain.GatewayErrorInterface)
)

type captureServiceMock struct{}

func (v captureServiceMock) CaptureTransactionAmount(request capture_domain.CaptureRequest) (*capture_domain.CaptureResponse, error_domain.GatewayErrorInterface) {
	return captureTransactionAmount(request)
}

func TestHandleCaptureRequest(t *testing.T) {
	expectedResponse := capture_domain.CaptureResponse{
		IsSuccess: true,
		Amount:    10,
		Currency:  "LKR",
	}

	captureTransactionAmount = func(request capture_domain.CaptureRequest) (*capture_domain.CaptureResponse, error_domain.GatewayErrorInterface) {
		return &expectedResponse, nil
	}

	capture_service.CaptureService = &captureServiceMock{}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)

	request := capture_domain.CaptureRequest{
		AuthId: "valid_string",
		Amount: 5,
	}

	b, err := json.Marshal(&request)
	if err != nil {
		t.Fail()
	}

	c.Request, err = http.NewRequest(http.MethodPatch, "", bytes.NewBuffer(b))
	if err != nil {
		t.Fail()
	}

	HandleCaptureRequest(c)
	var actualResponse capture_domain.CaptureResponse
	err = json.Unmarshal(response.Body.Bytes(), &actualResponse)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedResponse, actualResponse)
}

func TestHandleCaptureRequest_ErrorFromService(t *testing.T) {
	expectedError := error_domain.GatewayError{
		Code:  http.StatusUnprocessableEntity,
		Error: "error_from_service",
	}

	captureTransactionAmount = func(request capture_domain.CaptureRequest) (*capture_domain.CaptureResponse, error_domain.GatewayErrorInterface) {
		return nil, &expectedError
	}

	capture_service.CaptureService = &captureServiceMock{}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)

	request := capture_domain.CaptureRequest{
		AuthId: "valid_string",
		Amount: 5,
	}
	b, err := json.Marshal(&request)
	if err != nil {
		t.Fail()
	}

	c.Request, err = http.NewRequest(http.MethodPatch, "", bytes.NewBuffer(b))
	if err != nil {
		t.Fail()
	}

	HandleCaptureRequest(c)
	var actualError error_domain.GatewayError
	err = json.Unmarshal(response.Body.Bytes(), &actualError)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedError.ErrorMessage(), actualError.ErrorMessage())
}

func TestHandleCaptureRequest_InvalidBody(t *testing.T) {
	var err error
	expectedError := error_domain.GatewayError{
		Code:  http.StatusBadRequest,
		Error: "request body is invalid",
	}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)

	body := ioutil.NopCloser(strings.NewReader(`{
    													"rates": {
    													    "USD": "DON'T YOU PASS'",
															"GBP": "THIS IS SO WRONG"
														},
    													"base": "EUR",
   										 				"date": "2020-01-24"
														}`))

	c.Request, err = http.NewRequest(http.MethodPatch, "", body)
	if err != nil {
		t.Fail()
	}

	HandleCaptureRequest(c)
	var actualError error_domain.GatewayError
	err = json.Unmarshal(response.Body.Bytes(), &actualError)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedError.ErrorMessage(), actualError.ErrorMessage())
}
