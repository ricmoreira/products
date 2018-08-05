package errors

import (
	"products/models/response"
)

var (
	SERVICE_UNAVAILABLE string = "SERVICE_UNAVAILABLE"
	UNKNOWN_ERROR       string = "UNKNOWN_ERROR"
	DUPLICATED_ENTITY   string = "DUPLICATED_ENTITY"
	INVALID_REQUEST     string = "INVALID_REQUEST"
	EMPTY               string = "EMPTY"
	UNAUTHORIZED        string = "UNAUTHORIZED"
	NOT_FOUND           string = "NOT_FOUND"
)

var HttpErrorsMapper = map[string]int{
	SERVICE_UNAVAILABLE: 500,
	UNKNOWN_ERROR:       500,
	INVALID_REQUEST:     400,
	EMPTY:               400,
	DUPLICATED_ENTITY:   409,
	UNAUTHORIZED:        401,
	NOT_FOUND:           404,
}

var ResponseMessageErrorsMapper = map[string]string{
	SERVICE_UNAVAILABLE: "The service is currently unavailable",
	UNKNOWN_ERROR:       "Unknown server error",
	INVALID_REQUEST:     "Invalid request provided",
	EMPTY:               "Request with provided arguments resulted in an empty resource",
	DUPLICATED_ENTITY:   "Entity already exists",
	UNAUTHORIZED:        "User not found or invalid password",
	NOT_FOUND:           "This route does not exist",
}

// HandleErrorResponse returns a pointer to an ErrorResponse instance that matches App error response message protocol.
// If the httpCode is 0 we must use the default errorCode mapping association.
func HandleErrorResponse(errorCode string, errorsDetails []mresponse.ErrorDetail, customErrorMessage string) *mresponse.ErrorResponse {

	httpCode := HttpErrorsMapper[errorCode]
	if httpCode == 0 { // not found any matching Http code
		httpCode = 500
	}

	response := ResponseMessageErrorsMapper[errorCode]
	if response == "" { // not found any matching Response message
		response = "Unknown server error"
	}

	if customErrorMessage != "" { // replace custom error message if provided
		response = customErrorMessage
	}

	errorResponse := mresponse.ErrorResponse{
		HttpCode: httpCode,
		Code:     errorCode,
		Response: response,
	}

	if len(errorsDetails) != 0 {
		errorResponse.Errors = errorsDetails
	}

	return &errorResponse
}
