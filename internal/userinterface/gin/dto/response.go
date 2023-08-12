package dto

import "net/http"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type ErrorResponse struct {
	Status int     `json:"status"`
	Errors []Error `json:"error"`
}

// //
// http.StatusBadRequest
// //
const (
	BadRequest = iota
	RequestIDRequired
)

var (
	BadRequestResponse = ErrorResponse{Status: http.StatusBadRequest, Errors: []Error{
		{Code: BadRequest, Message: "bad request"},
	}}

	// X-Request-Id 헤더가 없을 때 반환되는 에러
	RequestIDRequiredResponse = ErrorResponse{Status: http.StatusBadRequest, Errors: []Error{
		{Code: RequestIDRequired, Message: "X-Request-Id header is required"},
	}}
)
