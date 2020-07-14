package void_service

import (
	"errors"
	"net/http"
	"payment-gateway-api/api/data_access"
	"payment-gateway-api/api/domain/error_domain"
	"payment-gateway-api/api/domain/void_domain"
)

type voidService struct{}

type voidServiceInterface interface {
	VoidTransaction(request void_domain.VoidRequest) (*void_domain.VoidResponse, error_domain.GatewayErrorInterface)
}

var (
	VoidService voidServiceInterface = &voidService{}
)

func (v *voidService) VoidTransaction(request void_domain.VoidRequest) (*void_domain.VoidResponse, error_domain.GatewayErrorInterface) {
	errs := request.ValidateFields()
	if len(errs) > 0 {
		return nil, error_domain.New(http.StatusUnprocessableEntity, errs...)
	}

	//check db for capture or refund operations
	ok, opsRecords, err := data_access.Db.GetAllOperationsByAuthID(request.AuthId)
	if err != nil {
		return nil, error_domain.New(http.StatusInternalServerError, errs...)
	}

	//if there are more than 1 operation with that unique auth ID, it means it must have
	//either be a void, a capture or a refund hence, the void cannot be executed.
	if ok && len(opsRecords) > 1 {
		return nil, error_domain.New(http.StatusUnprocessableEntity, errors.New("transaction is not in a state that allow cancellation"))
	}

	isPresent, authRecord, err := data_access.Db.GetAuthRecordByID(request.AuthId)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, error_domain.New(http.StatusNotFound, errors.New("authorisation transaction not found"))
		}
		return nil, error_domain.New(http.StatusInternalServerError, errors.New("unable to retrieve authorisation transaction"))
	}
	if !isPresent {
		return nil, error_domain.New(http.StatusOK, errors.New("transaction has already been cancelled"))
	}

	//otherwise we can soft delete the transaction by initialising the deletedAt field
	err = data_access.Db.SoftDeleteAuthRecordByID(request.AuthId)
	if err != nil {
		return nil, error_domain.New(http.StatusUnprocessableEntity, errors.New("internal issue preventing transaction to be cancelled"))
	}

	response := void_domain.VoidResponse{
		IsSuccess: true,
		Amount:    authRecord.AuthorisedAmount,
		Currency:  authRecord.Currency,
	}

	return &response, nil
}
