package helpers

import (
	"encoding/json"
	"net/http"
	"recibe_me/internal/models"
)

// Response returns a Response
func Response(responseWriter http.ResponseWriter, statusCode int, message string, errors interface{}, data interface{}) {
	response := models.ResponseModel{
		Message: message,
		Errors:  errors,
		Data:    data,
	}

	responseWriter.Header().Set("Content-type", "application/json; charset=utf-8")
	responseWriter.WriteHeader(statusCode)
	json.NewEncoder(responseWriter).Encode(response)
}
