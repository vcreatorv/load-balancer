package utils

import (
	"encoding/json"
	"errors"
	"lb/internal/models"
	"net/http"
)

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var errorToStatus = map[error]int{
	models.ErrNotFound:      http.StatusNotFound,
	models.ErrBadRequest:    http.StatusBadRequest,
	models.ErrForbidden:     http.StatusForbidden,
	models.ErrAlreadyExists: http.StatusConflict,
	models.ErrInternal:      http.StatusInternalServerError,
}

func GetErrorStatus(err error) int {
	var customErr *models.Error
	if errors.As(err, &customErr) {
		if status, ok := errorToStatus[customErr.ClientErr()]; ok {
			return status
		}
	}
	return http.StatusInternalServerError
}

func WriteError(w http.ResponseWriter, err error) {
	status := GetErrorStatus(err)
	apiError := APIError{
		Status:  status,
		Message: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiError)
}
