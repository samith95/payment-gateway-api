package refund_service

import (
	"errors"
	"log"
	"net/http"
	"payment-gateway-api/api/const/error_constant"
	"payment-gateway-api/api/data_access"
	"payment-gateway-api/api/data_access/database_model/auth"
	"payment-gateway-api/api/domain/common_validation"
	"payment-gateway-api/api/domain/error_domain"
	"payment-gateway-api/api/domain/refund_domain"
	"payment-gateway-api/api/services/common_service"
)

type refundService struct{}

type refundServiceInterface interface {
	RefundTransactionAmount(request refund_domain.RefundRequest) (*refund_domain.RefundResponse, error_domain.GatewayErrorInterface)
}

var (
	RefundService refundServiceInterface = &refundService{}
	operationName                        = "refund"
)

//RefundTransactionAmount refunds transaction amount of an already authorised and captured transaction by making sure the request and operations are valid
func (c *refundService) RefundTransactionAmount(request refund_domain.RefundRequest) (*refund_domain.RefundResponse, error_domain.GatewayErrorInterface) {
	//validate the refund operation
	authRecord, response, errInf := validateOperation(request)
	if errInf != nil {
		return response, errInf
	}

	//check that the amount is not greater than the amount that has been previously captured
	capturedAmount := authRecord.AuthorisedAmount - authRecord.AvailableAmount
	if request.Amount > capturedAmount {
		return nil, error_domain.New(http.StatusUnauthorized, errors.New(error_constant.RequestedAmountNotValid))
	}
	newAvailableAmount := authRecord.AvailableAmount + request.Amount

	//update available amount in db
	authRecord.AvailableAmount = newAvailableAmount
	err := data_access.Db.UpdateAvailableAmountByAuthID(authRecord.ID, newAvailableAmount, operationName)
	if err != nil {
		return nil, &error_domain.GatewayError{
			Code:  http.StatusInternalServerError,
			Error: error_constant.UpdateAvailableAmountFailure,
		}
	}

	return &refund_domain.RefundResponse{
		IsSuccess: true,
		Amount:    newAvailableAmount,
		Currency:  authRecord.Currency,
	}, nil
}

func validateOperation(request refund_domain.RefundRequest) (*auth.Auth, *refund_domain.RefundResponse, error_domain.GatewayErrorInterface) {
	errs := request.ValidateFields()
	if len(errs) > 0 {
		return nil, nil, error_domain.New(http.StatusUnprocessableEntity, errs...)
	}

	//check the requested operation is in execution during correct state
	isValid, err := common_service.CommonService.IsAuthorisedState(operationName, request.AuthId)
	if err != nil {
		log.Println(err.Error())
		return nil, nil, error_domain.New(http.StatusInternalServerError, errors.New(error_constant.UnableToCheckForInvalidState))
	}
	if !isValid {
		return nil, nil, error_domain.New(http.StatusUnprocessableEntity, errors.New(error_constant.TransactionStateInvalid))
	}

	isSoftDeleted, authRecord, err := data_access.Db.GetAuthRecordByID(request.AuthId)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == "record not found" {
			return nil, nil, error_domain.New(http.StatusNotFound, errors.New(error_constant.TransactionNotFound))
		}
		return nil, nil, error_domain.New(http.StatusInternalServerError, errors.New(error_constant.TransactionRetrievalFailure))
	}
	//check card number for refund failure reject
	isReject, err := data_access.Db.CheckRejectByCardNumber(operationName, authRecord.Number)
	if err != nil {
		log.Println(err.Error())
		return nil, nil, &error_domain.GatewayError{
			Code:  http.StatusInternalServerError,
			Error: error_constant.RejectRetrievalFailure,
		}
	}
	if isReject {
		return nil, nil, error_domain.New(http.StatusUnauthorized, errors.New(error_constant.RefundFailure))
	}

	//check transaction has been cancelled
	if !isSoftDeleted {
		return nil, nil, error_domain.New(http.StatusOK, errors.New(error_constant.CancelledTransaction))
	}
	//check expiration date, in case it was done at the end of the valid month
	if isValid := common_validation.IsExpiryDateValid(authRecord.ExpiryDate); !isValid {
		return nil, nil, error_domain.New(http.StatusUnauthorized, errors.New(error_constant.ExpiredCard))
	}
	return authRecord, nil, nil
}
