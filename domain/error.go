package domain

import "net/http"

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func BadRequestError(msg string) *Error {
	return &Error{
		http.StatusBadRequest,
		msg,
	}
}

func NotFoundError(msg string) *Error {
	return &Error{
		http.StatusNotFound,
		msg,
	}
}
