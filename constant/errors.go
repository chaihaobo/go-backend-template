package constant

import (
	"net/http"

	commonerr "github.com/chaihaobo/gocommon/error"
	"google.golang.org/grpc/codes"
)

//	Error constants; rule:  <projectID><moduleID><errorID>
//	projectID: 3digit
//	moduleID: 3digit
//	errorID: 4digit
//
// Successful always return "00000000"
var (
	Successful           = commonerr.ServiceError{Code: "0000000000", Message: "successful"}
	ErrSystemMalfunction = commonerr.ServiceError{Code: "9999999999", Message: "system malfunction"}
)

// ErrBadRequest bad request, user input error
var (
	ErrBadRequest   = commonerr.ServiceError{Code: "0010010001", Message: "bad request"}
	ErrUnauthorized = commonerr.ServiceError{Code: "0010010002", Message: "unauthorized"}
)

var (
	ErrHealthCheckFailed = commonerr.ServiceError{Code: "0010020001", Message: "health check failed"}
	ErrUserNotFound      = commonerr.ServiceError{Code: "0010020002", Message: "user not found"}
	ErrUserPasswordWrong = commonerr.ServiceError{Code: "0010020002", Message: "user password wrong"}
)

var ServiceErrorCode2HTTPStatus = map[string]int{
	Successful.Code:           http.StatusOK,
	ErrSystemMalfunction.Code: http.StatusInternalServerError,
	ErrBadRequest.Code:        http.StatusBadRequest,
	ErrHealthCheckFailed.Code: http.StatusServiceUnavailable,
	ErrUserNotFound.Code:      http.StatusBadRequest,
	ErrUnauthorized.Code:      http.StatusUnauthorized,
}

var ServiceErrorCodeToGRPCErrorCode = map[string]codes.Code{
	Successful.Code: codes.OK,
}
