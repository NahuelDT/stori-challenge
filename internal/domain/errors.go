package domain

import "errors"

var (
	ErrInvalidTransaction   = errors.New("invalid transaction data")
	ErrInvalidAmount        = errors.New("invalid transaction amount")
	ErrInvalidDate          = errors.New("invalid transaction date")
	ErrInvalidTransactionID = errors.New("invalid transaction ID")
	ErrFileNotFound         = errors.New("transaction file not found")
	ErrInvalidFileFormat    = errors.New("invalid file format")
	ErrEmailDeliveryFailed  = errors.New("email delivery failed")
	ErrDatabaseConnection   = errors.New("database connection failed")
)
