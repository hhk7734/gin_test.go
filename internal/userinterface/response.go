package userinterface

type RequestError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type ErrorResponse struct {
	Error RequestError `json:"error"`
}

func (r ErrorResponse) Message(m string) ErrorResponse {
	return ErrorResponse{Error: RequestError{Code: r.Error.Code, Message: m}}
}

func (r ErrorResponse) Code(c int) ErrorResponse {
	return ErrorResponse{Error: RequestError{Code: c, Message: r.Error.Message}}
}

// //
// http.StatusBadRequest
// //
const (
	_ int = iota
	RequestIDRequired
)

var (
	// X-Request-Id 헤더가 없을 때 반환되는 에러
	RequestIDRequiredResponse = ErrorResponse{Error: RequestError{Code: RequestIDRequired, Message: "X-Request-Id header is required"}}
)
