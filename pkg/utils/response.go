package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"merge-queue/internal/models"
)

// ResponseHelper provides utility functions for HTTP responses.
type ResponseHelper struct{}

// NewResponseHelper creates a new ResponseHelper instance.
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// SendJSON sends a JSON response.
func (rh *ResponseHelper) SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// SendError sends an error response.
func (rh *ResponseHelper) SendError(w http.ResponseWriter, statusCode int, message string) {
	response := models.APIResponse{
		Success:   false,
		Error:     message,
		Timestamp: time.Now(),
	}
	rh.SendJSON(w, statusCode, response)
}

// SendErrorWithCode sends an error response with a specific error code.
func (rh *ResponseHelper) SendErrorWithCode(w http.ResponseWriter, statusCode int, code, message, details string) {
	errorResp := models.ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}

	response := models.APIResponse{
		Success:   false,
		Error:     message,
		Data:      errorResp,
		Timestamp: time.Now(),
	}

	rh.SendJSON(w, statusCode, response)
}

// SendSuccess sends a success response.
func (rh *ResponseHelper) SendSuccess(w http.ResponseWriter, data interface{}) {
	response := models.APIResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now(),
	}
	rh.SendJSON(w, http.StatusOK, response)
}

// SendSuccessWithMeta sends a success response with metadata.
func (rh *ResponseHelper) SendSuccessWithMeta(w http.ResponseWriter, data interface{}, meta interface{}) {
	response := models.APIResponse{
		Success:   true,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now(),
	}
	rh.SendJSON(w, http.StatusOK, response)
}

// SendCreated sends a 201 Created response.
func (rh *ResponseHelper) SendCreated(w http.ResponseWriter, data interface{}) {
	response := models.APIResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now(),
	}
	rh.SendJSON(w, http.StatusCreated, response)
}

// SendNoContent sends a 204 No Content response.
func (rh *ResponseHelper) SendNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// SendPaginated sends a paginated response with metadata.
func (rh *ResponseHelper) SendPaginated(w http.ResponseWriter, data interface{}, page, perPage, total int) {
	totalPages := (total + perPage - 1) / perPage // Ceiling division.

	meta := models.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}

	rh.SendSuccessWithMeta(w, data, meta)
}
