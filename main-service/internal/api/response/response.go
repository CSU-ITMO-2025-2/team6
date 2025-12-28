package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"main-service/internal/model"
)

type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func OK(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, data)
}

func Error(w http.ResponseWriter, httpCode int, code string, message string) {
	JSON(w, httpCode, ErrorBody{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func HandleDomainError(w http.ResponseWriter, err error) {
	httpCode, code, msg := MapDomainError(err)
	Error(w, httpCode, code, msg)
}

func MapDomainError(err error) (int, string, string) {
	//var dErr domain_error.DomainError
	//if errors.As(err, &dErr) {
	//	return dErr.HTTPStatus(), dErr.Code(), dErr.Message()
	//}

	switch {
	case errors.Is(err, model.ErrObjectNotFound):
		return http.StatusNotFound, "NOT_FOUND", err.Error()

	case errors.Is(err, model.ErrObjectAlreadyExists):
		return http.StatusConflict, "ALREADY_EXISTS", err.Error()

	case errors.Is(err, model.ErrInvalidInput):
		return http.StatusBadRequest, "INVALID_INPUT", err.Error()

	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", err.Error()
	}
}
