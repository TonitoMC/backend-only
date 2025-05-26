package handlers

// For this project I've implemented a custom HTTP error handler, this works by having the
// HTTP handlers return errors and mapping here to specific responses. The errors are declared
// internal/errors/errors.go, and throughout the domain layer they're wrapped with context as
// these generic errors for not found, conflict, etc. The naming of 'domain errors' was a bit of
// an oversight as we also do things like transform URL parameters to integers inside of the
// handlers which could return these, I'd consider adding a global errors package or taking a
// different approach for future projects

import (
	"errors"
	"net/http"
	"strings"

	domainErrors "series-tracker/internal/errors"
	"series-tracker/internal/models"

	"github.com/labstack/echo/v4"
)

// Custom error handler, takes an error and maps to a response
func CustomHTTPErrorHandler(err error, c echo.Context) {
	// Status - HTTP Status, Code - message_written_like_this, Message - Message written like this
	status, code, message := mapError(err)

	// Builds response with code & message
	errResponse := models.ErrorResponse{
		Code:    code,
		Message: message,
	}

	c.JSON(status, errResponse)
}

// MapError takes care of all the mappings of the errors declared in /errors/errors.go, each specific
// error is wrapped with more information which we don't want to give off to anyone using the API. What
// we're doing here is just exposing the most basic of information that also gives off a signal on what
// went wrong, any other errors are mapped to a generic 500 internal.
func mapError(err error) (int, string, string) {
	if errors.Is(err, domainErrors.ErrSeriesNotFound) {
		return http.StatusNotFound, "series_not_found", capitalize(domainErrors.ErrSeriesNotFound.Error())
	}

	if errors.Is(err, domainErrors.ErrInvalidInput) {
		return http.StatusBadRequest, "invalid_input", capitalize(domainErrors.ErrInvalidInput.Error())
	}

	if errors.Is(err, domainErrors.ErrSeriesConflict) {
		return http.StatusConflict, "series_conflict", capitalize(domainErrors.ErrSeriesConflict.Error())
	}

	return http.StatusInternalServerError, "internal_server_error", "Internal server error"
}

// Utility function for prettier messages
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
