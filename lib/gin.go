package lib

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Helpers for GIN
type ValidationDetail struct {
	Rule    string `json:"rules"`
	Message string `json:"error_message"`
}

type ValidationErrorResponse struct {
	Error  bool                        `json:"error"`
	Errors map[string]ValidationDetail `json:"errors_detail"`
}

func IsValidRequest(context *gin.Context, validStructure any) bool {

	if err := context.ShouldBindJSON(&validStructure); err != nil {

		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			// Handle other types of errors (e.g., parsing errors)
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return false
		}

		response := ValidationErrorResponse{
			Error:  true,
			Errors: make(map[string]ValidationDetail),
		}

		for _, validationError := range validationErrors {

			validStructName := strings.ToLower(GetStructName(validStructure))
			namespace := strings.ToLower(validationError.StructNamespace())
			field := strings.Replace(namespace, validStructName+".", "", 1)
			rule := strings.ToLower(validationError.Tag())

			response.Errors[field] = ValidationDetail{
				Rule:    rule,
				Message: validationError.Error(),
			}
		}

		context.JSON(http.StatusBadRequest, response)
		return false
	}

	return true
}
