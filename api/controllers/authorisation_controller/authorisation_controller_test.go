package authorisation_controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"payment-gateway-api/api/domain/auth_domain"
	"payment-gateway-api/api/domain/error_domain"
	"payment-gateway-api/api/services/authorisation_service"
	"strings"
	"testing"
)

var (
	authoriseTransactionFunc func(auth_domain.AuthRequest) (*auth_domain.AuthResponse, error_domain.GatewayErrorInterface)
)

type authoriseServiceMock struct{}

func (a *authoriseServiceMock) GetAllRecords() (string, error_domain.GatewayErrorInterface) {
	panic("implement me")
}

func (a *authoriseServiceMock) AuthoriseTransaction(request auth_domain.AuthRequest) (*auth_domain.AuthResponse, error_domain.GatewayErrorInterface) {
	return authoriseTransactionFunc(request)
}

func TestHandleAuthorisationRequestSuccess(t *testing.T) {
	expectedResponse := auth_domain.AuthResponse{
		AuthID:    "valid_auth_id",
		IsSuccess: true,
		Amount:    10,
		Currency:  "GBP",
	}

	authoriseTransactionFunc = func(request auth_domain.AuthRequest) (*auth_domain.AuthResponse, error_domain.GatewayErrorInterface) {
		return &expectedResponse, nil
	}

	authorisation_service.AuthorisationService = &authoriseServiceMock{}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)

	cardDetails := auth_domain.CardDetails{
		Number:     "4929907390318794",
		ExpiryDate: "12-3500",
		Cvv:        "123",
	}
	request := auth_domain.AuthRequest{
		CardDetails: cardDetails,
		Amount:      10000,
		Currency:    "GBP",
	}

	b, err := json.Marshal(&request)
	if err != nil {
		t.Fail()
	}

	c.Request, err = http.NewRequest(http.MethodPost, "", bytes.NewBuffer(b))
	if err != nil {
		t.Fail()
	}

	HandleAuthorisationRequest(c)
	var actualResponse auth_domain.AuthResponse
	err = json.Unmarshal(response.Body.Bytes(), &actualResponse)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedResponse, actualResponse)
}

func TestHandleAuthorisationRequestErrorFromService(t *testing.T) {
	expectedError := error_domain.GatewayError{
		Code:  http.StatusUnprocessableEntity,
		Error: "error_from_service",
	}

	authoriseTransactionFunc = func(request auth_domain.AuthRequest) (*auth_domain.AuthResponse, error_domain.GatewayErrorInterface) {
		return nil, &expectedError
	}

	authorisation_service.AuthorisationService = &authoriseServiceMock{}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)

	cardDetails := auth_domain.CardDetails{
		Number:     "4929907390318794",
		ExpiryDate: "12-3500",
		Cvv:        "123",
	}
	request := auth_domain.AuthRequest{
		CardDetails: cardDetails,
		Amount:      10000,
		Currency:    "GBP",
	}

	b, err := json.Marshal(&request)
	if err != nil {
		t.Fail()
	}

	c.Request, err = http.NewRequest(http.MethodPost, "", bytes.NewBuffer(b))
	if err != nil {
		t.Fail()
	}

	HandleAuthorisationRequest(c)
	var actualError error_domain.GatewayError
	err = json.Unmarshal(response.Body.Bytes(), &actualError)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedError, actualError)
}

func TestHandleAuthorisationRequestInvalidBody(t *testing.T) {
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

	c.Request, err = http.NewRequest(http.MethodPost, "", body)
	if err != nil {
		t.Fail()
	}

	HandleAuthorisationRequest(c)
	var actualError error_domain.GatewayError
	err = json.Unmarshal(response.Body.Bytes(), &actualError)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedError, actualError)
}
