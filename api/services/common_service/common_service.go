package common_service

import (
	"errors"
	"log"
	"payment-gateway-api/api/const/error_constant"
	"payment-gateway-api/api/data_access"
)

type commonService struct{}

type commonServiceInterface interface {
	IsAuthorisedState(string, string) (bool, error)
}

var (
	CommonService commonServiceInterface = &commonService{}
)

//IsAuthorisedState will check whether the operation required by the client are
//authorised in relation to the auth transaction lifecycle diagram
func (c *commonService) IsAuthorisedState(operationName, id string) (bool, error) {
	var invalidPreviousState string

	switch operationName {
	case "void":
		invalidPreviousState = "capture"
	case "capture":
		invalidPreviousState = "refund"
	case "refund":
		return true, nil
	default:
		return false, errors.New(error_constant.OperationNameInvalid)
	}

	//check whether previous state that are invalid for the current operation are present in db
	isPresent, _, err := data_access.Db.GetOperationByAuthIDAndOperationName(id, invalidPreviousState)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	//if they are, return error stating invalid state
	if isPresent {
		return false, nil
	}

	return true, nil
}
