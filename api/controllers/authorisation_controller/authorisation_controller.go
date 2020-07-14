package authorisation_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"payment-gateway-api/api/domain/auth_domain"
	"payment-gateway-api/api/domain/error_domain"
	"payment-gateway-api/api/services/authorisation_service"
)

func HandleAuthorisationRequest(c *gin.Context) {
	request := auth_domain.AuthRequest{}

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, error_domain.GatewayError{
			Code:  http.StatusBadRequest,
			Error: "request body is invalid",
		})
		return
	}

	result, apiError := authorisation_service.AuthorisationService.AuthoriseTransaction(request)
	if apiError != nil {
		c.JSON(apiError.Status(), apiError)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func HandleGetAuthsRecords(c *gin.Context) {
	result, apiError := authorisation_service.AuthorisationService.GetAllRecords()
	if apiError != nil {
		c.JSON(apiError.Status(), apiError)
		return
	}
	c.JSON(http.StatusOK, result)
}
