package xerror

import (
	"errors"
	"net/http"
)

// ErrorSession represents authentication or session-related issues
type ErrorSession struct {
	Message string
}

func (e *ErrorSession) Error() string { return e.Message }

// ErrorNotFound represents missing resources
type ErrorNotFound struct {
	Message string
}

func (e *ErrorNotFound) Error() string { return e.Message }

// ErrorPermission represents authorization (RBAC) issues
type ErrorPermission struct {
	Message string
}

func (e *ErrorPermission) Error() string { return e.Message }

// ErrorPermission represents authorization (RBAC) issues
type ErrorBadRequest struct {
	Message string
}

func (e *ErrorBadRequest) Error() string { return e.Message }

// DefineStatusCode maps custom error types to HTTP Status Codes
func DefineStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	// Use errors.As to detect wrapped errors or specific types
	var errSession *ErrorSession
	if errors.As(err, &errSession) {
		return http.StatusUnauthorized
	}

	var errPermission *ErrorPermission
	if errors.As(err, &errPermission) {
		return http.StatusForbidden
	}

	var errNotFound *ErrorNotFound
	if errors.As(err, &errNotFound) {
		return http.StatusNotFound
	}

	// Fallback for everything else
	return http.StatusInternalServerError
}
