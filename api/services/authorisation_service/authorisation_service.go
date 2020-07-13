package authorisation_service

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	dal "payment-gateway-api/api/data_access"
	"payment-gateway-api/api/data_access/database_model"
	"payment-gateway-api/api/domain/gateway_domain/auth_domain"
	"payment-gateway-api/api/domain/gateway_domain/error_domain"
	"time"
)

type authorisationService struct{}

type authorisationServiceInterface interface {
	AuthorisePayment(input auth_domain.AuthRequest) (*auth_domain.AuthResponse, error_domain.GatewayErrorInterface)
	GetAllRecords() (string, error_domain.GatewayErrorInterface)
}

var (
	AuthorisationService authorisationServiceInterface = &authorisationService{}
)

func (a *authorisationService) AuthorisePayment(request auth_domain.AuthRequest) (*auth_domain.AuthResponse, error_domain.GatewayErrorInterface) {
	errs := request.ValidateFields()
	if len(errs) > 0 {
		return nil, error_domain.New(http.StatusBadRequest, errs...)
	}

	//generate uniqueID
	authId := uuid.New().String()

	record := database_model.Auth{
		ID:               authId,
		Number:           request.CardDetails.Number,
		ExpiryDate:       request.CardDetails.ExpiryDate,
		AuthorisedAmount: request.Amount,
		AvailableAmount:  request.Amount,
		Currency:         request.Currency,
		AuditData: database_model.AuditData{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	err := dal.Db.CreateAuthRecord(&record)
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
	err := dal.Db.Init()
	if err != nil {
		return "", error_domain.New(http.StatusBadRequest, err)
	}
	var records []database_model.Auth

	records, err = dal.Db.GetAllRecords()
	if err != nil {
		return "", error_domain.New(http.StatusBadRequest, err)
	}

	recordsStr := fmt.Sprintf("%+v\n", records)

	return recordsStr, nil
}
