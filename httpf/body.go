package httpf

import (
	"net/http"

	"github.com/Prastiwar/Go-flow/datas"
)

var (
	_ BodyUnmarshaler = &bodyUnmarshaler{}
)

// BodyUnmarshaler is an interface that defines a method to unmarshal the body of an HTTP response into a value of any type.
// The implementation of this interface should handle the cases where the response status code indicates an error and return an appropriate error value.
type BodyUnmarshaler interface {
	Unmarshal(r *http.Response, v any) error
}

type bodyUnmarshaler struct {
	errorHandler func(r *http.Response, u datas.ReaderUnmarshaler) error
	unmarshaler  datas.ReaderUnmarshaler
}

// NewBodyUnmarshaler returns an implementation for BodyUnmarshaler which unmarshals response body and supports error handling.
// Passed unmarshaler will be used to unmarshal both, actual body value or error value from body response.
// If IsErrorStatus will return true during Unmarshal, it will call the errorHandler to transform error or handle it.
func NewBodyUnmarshaler(u datas.ReaderUnmarshaler, errorHandler func(r *http.Response, u datas.ReaderUnmarshaler) error) BodyUnmarshaler {
	return &bodyUnmarshaler{
		errorHandler: errorHandler,
		unmarshaler:  u,
	}
}

// NewBodyUnmarshalerWithError returns an implementation for BodyUnmarshaler which unmarshals response body and supports http error unmarshaling.
// Passed unmarshaler will be used to unmarshal both, actual body value or error value from body response.
// errorStruct must be a pointer to error struct which indicates structure of the error returned from API in body response.
// If IsErrorStatus will return true during Unmarshal, it will copy errorStruct and unmarshal error response into copied value and return it.
// If custom error handling e.g for 404 codes is needed, you should use NewBodyUnmarshaler.
func NewBodyUnmarshalerWithError(u datas.ReaderUnmarshaler, errorStruct error) BodyUnmarshaler {
	return NewBodyUnmarshaler(u, func(r *http.Response, u datas.ReaderUnmarshaler) error {
		var httpError = errorStruct
		if err := u.UnmarshalFrom(r.Body, &httpError); err != nil {
			return err
		}
		return httpError
	})
}

func (u *bodyUnmarshaler) Unmarshal(r *http.Response, v any) error {
	defer r.Body.Close()

	if IsErrorStatus(r.StatusCode) {
		return u.errorHandler(r, u.unmarshaler)
	}

	return u.unmarshaler.UnmarshalFrom(r.Body, v)
}
