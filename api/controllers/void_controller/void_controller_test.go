package void_controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"payment-gateway-api/api/domain/error_domain"
	"payment-gateway-api/api/domain/void_domain"
	"payment-gateway-api/api/services/void_service"
	"strings"
	"testing"
)

var (
	voidTransaction func(void_domain.VoidRequest) (*void_domain.VoidResponse, error_domain.GatewayErrorInterface)
)

type voidServiceMock struct{}

func (v voidServiceMock) VoidTransaction(request void_domain.VoidRequest) (*void_domain.VoidResponse, error_domain.GatewayErrorInterface) {
	return voidTransaction(request)
}

func TestHandleVoidRequest(t *testing.T) {
	expectedResponse := void_domain.VoidResponse{
		IsSuccess: true,
		Amount:    10,
		Currency:  "LKR",
	}

	voidTransaction = func(request void_domain.VoidRequest) (response *void_domain.VoidResponse, errorInterface error_domain.GatewayErrorInterface) {
		return &expectedResponse, nil
	}

	void_service.VoidService = &voidServiceMock{}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)

	request := void_domain.VoidRequest{AuthId: "valid_string"}

	b, err := json.Marshal(&request)
	if err != nil {
		t.Fail()
	}

	c.Request, err = http.NewRequest(http.MethodPatch, "", bytes.NewBuffer(b))
	if err != nil {
		t.Fail()
	}

	HandleVoidRequest(c)
	var actualResponse void_domain.VoidResponse
	err = json.Unmarshal(response.Body.Bytes(), &actualResponse)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedResponse, actualResponse)
}

func TestHandleVoidRequest_ErrorFromService(t *testing.T) {
	expectedError := error_domain.GatewayError{
		Code:  http.StatusUnprocessableEntity,
		Error: "error_from_service",
	}

	voidTransaction = func(request void_domain.VoidRequest) (response *void_domain.VoidResponse, errorInterface error_domain.GatewayErrorInterface) {
		return nil, &expectedError
	}

	void_service.VoidService = &voidServiceMock{}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)

	request := void_domain.VoidRequest{AuthId: "valid_string"}

	b, err := json.Marshal(&request)
	if err != nil {
		t.Fail()
	}

	c.Request, err = http.NewRequest(http.MethodPatch, "", bytes.NewBuffer(b))
	if err != nil {
		t.Fail()
	}

	HandleVoidRequest(c)
	var actualError error_domain.GatewayError
	err = json.Unmarshal(response.Body.Bytes(), &actualError)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedError, actualError)
}

func TestHandleVoidRequest_InvalidBody(t *testing.T) {
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

	HandleVoidRequest(c)
	var actualError error_domain.GatewayError
	err = json.Unmarshal(response.Body.Bytes(), &actualError)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedError, actualError)
}
