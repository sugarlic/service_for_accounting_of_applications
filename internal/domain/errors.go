package domain

import "errors"

var (
	ErrApplicationNotFound = errors.New("application not found")
	ErrInvalidApplication  = errors.New("invalid application")
	ErrInvalidStatus       = errors.New("invalid application status")
	ErrNameRequired        = errors.New("name is required")
	ErrPhoneRequired       = errors.New("phone is required")
)
