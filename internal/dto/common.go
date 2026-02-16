package dto

import "errors"

const (
	// Failed Messages
	MESSAGE_FAILED_PROSES_REQUEST      = "failed to process request"
	MESSAGE_FAILED_TOKEN_NOT_FOUND     = "token not found"
	MESSAGE_FAILED_TOKEN_NOT_VALID     = "token not valid"
	MESSAGE_FAILED_DENIED_ACCESS       = "denied access"
	MESSAGE_FAILED_PARSE_TIME          = "failed to parse time"
	MESSAGE_FAILED_GET_DATA_FROM_BODY  = "failed to get data from body"
	MESSAGE_FAILED_GET_CALLBACK_TRIPAY = "failed to get callback from tripay"

	// Success Messages
	MESSAGE_SUCCESS_GET_CALLBACK_TRIPAY = "success get callback from tripay"

	// General Messages
	PESAN_DILUAR_MASA_REGISTRASI = "request made outside of allowed time frame"
)

var (
	ErrRoleNotAllowed = errors.New("role not allowed")
)
