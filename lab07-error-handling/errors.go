package main

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

func (e *APIError) Send(w http.ResponseWriter) {
	w.WriteHeader(e.StatusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: *e})
}

func NewBadRequestError(message string) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    message,
	}
}

func NewNotFoundError(resource string) *APIError {
	return &APIError{
		StatusCode: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    resource + " not found",
	}
}

func NewConflictError(message string) *APIError {
	return &APIError{
		StatusCode: http.StatusConflict,
		Code:       "CONFLICT",
		Message:    message,
	}
}

func NewInternalError() *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_ERROR",
		Message:    "An internal server error occurred",
	}
}
