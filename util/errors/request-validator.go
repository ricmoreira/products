package errors

import (
	"products/models/response"

	"github.com/asaskevich/govalidator"
)

func ValidateRequest(obj interface{}) *mresponse.ErrorResponse {

	_, er := govalidator.ValidateStruct(obj)

	if er != nil {
		e := mresponse.ErrorResponse{}
		e.HttpCode = 400
		e.Code = INVALID_REQUEST
		e.Response = ResponseMessageErrorsMapper[INVALID_REQUEST]

		details := []mresponse.ErrorDetail{}
		for k, v := range govalidator.ErrorsByField(er) {
			errorDetail := mresponse.ErrorDetail{
				Property: k,
				Message:  v,
			}
			details = append(details, errorDetail)
		}

		e.Errors = details

		return &e
	}

	return nil
}
