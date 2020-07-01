package errors

import (
	"math"
)

type ErrorCode uint16

const NoCode ErrorCode = math.MaxUint16
const BadRequest ErrorCode = 400
const NotFound ErrorCode = 404

func cause(err error) error {
	type causer interface {
		Cause() error
	}

	if cause, ok := err.(causer); ok {
		return cause.Cause()
	}
	return nil
}

func code(err error) ErrorCode {

	if coder, ok := err.(interface{ Code() ErrorCode }); ok {
		return coder.Code()
	}
	return NoCode
}
func GetCode(err error) (ErrorCode, string) {

	for ; err != nil; err = cause(err) {
		if coder, ok := err.(interface{ Code() ErrorCode }); ok {
			return coder.Code(), err.Error()
		}
	}
	return NoCode, ""

}
