package tools

import (
	"encoding/json"
	"errors"
	"net/http"

	commonerr "github.com/chaihaobo/gocommon/error"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"

	"gitlab.seakoi.net/engineer/backend/be-template/constant"
	govalidator "gitlab.seakoi.net/engineer/backend/be-template/resource/validator"
)

type responseBody struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Data    any    `json:"data"`
}

func newResponseBody(code, message string, data any) responseBody {
	return responseBody{code, message, data}
}

func normalResponseBody(data any) responseBody {
	return newResponseBody(constant.Successful.Code, constant.Successful.Message, data)
}

// HTTPWrite write a normal json response
func HTTPWrite(writer http.ResponseWriter, data any) {
	writer.WriteHeader(http.StatusOK)
	appendGenericHeader(writer)
	if data == nil {
		return
	}
	writer.Write(lo.Must(json.Marshal(normalResponseBody(data))))
}

func appendGenericHeader(writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json;charset=utf=8")
}

// HTTPWriteErr write a error json body
func HTTPWriteErr(writer http.ResponseWriter, err error) {
	var (
		serviceErr = constant.ErrSystemMalfunction
	)

	switch err := err.(type) {
	case commonerr.ServiceError:
		serviceErr = err
	case validator.ValidationErrors:
		serviceErr = commonerr.ServiceError{
			Code: constant.ErrBadRequest.Code,
			Message: lo.Entries[string, string](err.
				Translate(govalidator.Translator))[0].Value,
		}
	default:
		serviceErr = constant.ErrSystemMalfunction
	}
	_ = errors.As(err, &serviceErr)
	actualHTTPStatus, ok := constant.ServiceErrorCode2HTTPStatus[serviceErr.Code]
	if !ok {
		actualHTTPStatus = http.StatusInternalServerError
	}

	for key, value := range serviceErr.Attributes {
		writer.Header().Add(key, value)
	}
	writer.WriteHeader(actualHTTPStatus)
	writer.Write(lo.Must(json.Marshal(newResponseBody(serviceErr.Code, serviceErr.Message, nil))))
}
