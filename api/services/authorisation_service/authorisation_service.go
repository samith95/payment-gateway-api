package authorisation_service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"payment-gateway-api/api/const/error_constant"
	dal "payment-gateway-api/api/data_access"
	"payment-gateway-api/api/data_access/database_model/auth"
	"payment-gateway-api/api/domain/auth_domain"
	"payment-gateway-api/api/domain/error_domain"
	"time"
)

type authorisationService struct{}

type authorisationServiceInterface interface {
	AuthoriseTransaction(auth_domain.AuthRequest) (*auth_domain.AuthResponse, error_domain.GatewayErrorInterface)
	GetAllRecords() (string, error_domain.GatewayErrorInterface)
}

var (
	AuthorisationService authorisationServiceInterface = &authorisationService{}
)

func (a *authorisationService) AuthoriseTransaction(request auth_domain.AuthRequest) (*auth_domain.AuthResponse, error_domain.GatewayErrorInterface) {
	errs := request.ValidateFields()
	if len(errs) > 0 {
		return nil, error_domain.New(http.StatusBadRequest, errs...)
	}

	isReject, err := dal.Db.CheckRejectByCardNumber("authorisation", request.CardDetails.Number)
	if err != nil {
		return nil, &error_domain.GatewayError{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		}
	}
	if isReject {
		return nil, error_domain.New(http.StatusUnauthorized, errors.New(error_constant.AuthorisationFailure))
	}

	//generate uniqueID
	authId := uuid.New().String()

	record := auth.Auth{
		ID:               authId,
		Number:           request.CardDetails.Number,
		ExpiryDate:       request.CardDetails.ExpiryDate,
		AuthorisedAmount: request.Amount,
		AvailableAmount:  request.Amount,
		Currency:         request.Currency,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		DeletedAt:        time.Time{},
	}

	err = dal.Db.InsertAuthRecord(&record)
	if err != nil {
		return nil, &error_domain.GatewayError{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		}
	}

	response := auth_domain.AuthResponse{
		AuthID:    authId,
		IsSuccess: true,
		Amount:    request.Amount,
		Currency:  request.Currency,
	}

	return &response, nil
}

func (a *authorisationService) GetAllRecords() (string, error_domain.GatewayErrorInterface) {
	var records []auth.Auth

	records, err := dal.Db.GetAllAuthRecords()
	if err != nil {
		return "", error_domain.New(http.StatusBadRequest, err)
	}

	recordsStr := fmt.Sprintf("%+v\n", records)

	return recordsStr, nil
}
