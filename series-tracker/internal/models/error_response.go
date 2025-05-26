package models

// ErrorResponse defines the structure of API error responses.
type ErrorResponse struct {
	Code    string `json:"code" example:"lowercase_and_underscores"`    // Machine-friendly error message for handling in front-end
	Message string `json:"message" example:"Capitalization and spaces"` // Human-friendly error message for the user
}
