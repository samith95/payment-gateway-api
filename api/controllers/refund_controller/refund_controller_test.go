package refund_controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"payment-gateway-api/api/domain/error_domain"
	"payment-gateway-api/api/domain/refund_domain"
	"payment-gateway-api/api/services/refund_service"
	"strings"
	"testing"
)

var (
	refundTransactionAmount func(request refund_domain.RefundRequest) (*refund_domain.RefundResponse, error_domain.GatewayErrorInterface)
)

type refundServiceMock struct{}

func (v refundServiceMock) RefundTransactionAmount(request refund_domain.RefundRequest) (*refund_domain.RefundResponse, error_domain.GatewayErrorInterface) {
	return refundTransactionAmount(request)
}

func TestHandleRefundRequest(t *testing.T) {
	expectedResponse := refund_domain.RefundResponse{
		IsSuccess: true,
		Amount:    10,
		Currency:  "LKR",
	}

	refundTransactionAmount = func(request refund_domain.RefundRequest) (*refund_domain.RefundResponse, error_domain.GatewayErrorInterface) {
		return &expectedResponse, nil
	}

	refund_service.RefundService = &refundServiceMock{}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)

	request := refund_domain.RefundRequest{
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

	HandleRefundRequest(c)
	var actualResponse refund_domain.RefundResponse
	err = json.Unmarshal(response.Body.Bytes(), &actualResponse)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedResponse, actualResponse)
}

func TestHandleRefundRequest_ErrorFromService(t *testing.T) {
	expectedError := error_domain.GatewayError{
		Code:  http.StatusUnprocessableEntity,
		Error: "error_from_service",
	}

	refundTransactionAmount = func(request refund_domain.RefundRequest) (*refund_domain.RefundResponse, error_domain.GatewayErrorInterface) {
		return nil, &expectedError
	}

	refund_service.RefundService = &refundServiceMock{}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)

	request := refund_domain.RefundRequest{
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

	HandleRefundRequest(c)
	var actualError error_domain.GatewayError
	err = json.Unmarshal(response.Body.Bytes(), &actualError)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedError.ErrorMessage(), actualError.ErrorMessage())
}

func TestHandleRefundRequest_InvalidBody(t *testing.T) {
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

	HandleRefundRequest(c)
	var actualError error_domain.GatewayError
	err = json.Unmarshal(response.Body.Bytes(), &actualError)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedError.ErrorMessage(), actualError.ErrorMessage())
}
