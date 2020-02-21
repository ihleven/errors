package errors

import (
	"math"
)

type ErrorCode int

const NoCode ErrorCode = math.MaxUint16
const BadRequest ErrorCode = 400
const NotFound ErrorCode = 404

type coder interface {
	Code() int
}

func GetCauseCode(err error) int {

	cause := Cause(err)

	if errWithCode, ok := cause.(coder); ok {
		return errWithCode.Code()
	}

	return int(NoCode)
}

func Code(err error) int {

	type causer interface {
		Cause() error
	}

	for err != nil {
		errWithCode, ok := err.(coder)
		if ok {
			return int(errWithCode.Code())
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return int(NoCode)
}
