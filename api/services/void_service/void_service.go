package void_service

import (
	"errors"
	"net/http"
	"payment-gateway-api/api/const/error_constant"
	"payment-gateway-api/api/data_access"
	"payment-gateway-api/api/domain/error_domain"
	"payment-gateway-api/api/domain/void_domain"
	"payment-gateway-api/api/services/common_service"
)

type voidService struct{}

type voidServiceInterface interface {
	VoidTransaction(request void_domain.VoidRequest) (*void_domain.VoidResponse, error_domain.GatewayErrorInterface)
}

var (
	VoidService   voidServiceInterface = &voidService{}
	operationName                      = "void"
)

//VoidTransaction cancels a transaction after being authorised by making sure the request and operations are valid
func (v *voidService) VoidTransaction(request void_domain.VoidRequest) (*void_domain.VoidResponse, error_domain.GatewayErrorInterface) {
	errs := request.ValidateFields()
	if len(errs) > 0 {
		return nil, error_domain.New(http.StatusUnprocessableEntity, errs...)
	}

	//check operation can be executed according to state
	isValid, err := common_service.CommonService.IsAuthorisedState(operationName, request.AuthId)
	if err != nil {
		return nil, error_domain.New(http.StatusInternalServerError, errors.New(error_constant.UnableToCheckForInvalidState))
	}
	if !isValid {
		return nil, error_domain.New(http.StatusUnprocessableEntity, errors.New(error_constant.TransactionStateInvalid))
	}

	isSoftDeleted, authRecord, err := data_access.Db.GetAuthRecordByID(request.AuthId)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, error_domain.New(http.StatusNotFound, errors.New(error_constant.TransactionNotFound))
		}
		return nil, error_domain.New(http.StatusInternalServerError, errors.New(error_constant.TransactionRetrievalFailure))
	}
	if !isSoftDeleted {
		return nil, error_domain.New(http.StatusOK, errors.New("transaction has already been cancelled"))
	}

	//otherwise we can soft delete the transaction by initialising the deletedAt field
	err = data_access.Db.SoftDeleteAuthRecordByID(request.AuthId)
	if err != nil {
		return nil, error_domain.New(http.StatusUnprocessableEntity, errors.New(error_constant.UnableToVoidTransaction))
	}

	response := void_domain.VoidResponse{
		IsSuccess: true,
		Amount:    authRecord.AuthorisedAmount,
		Currency:  authRecord.Currency,
	}

	return &response, nil
}
