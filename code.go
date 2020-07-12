package errors

import (
	"math"
)

// ErrorCode is a code that can be attached to an error and is passed up the stack.
type ErrorCode int

// NoCode is the error code of errors with no code explicitly attached.
// It has a value of math.MaxUint16. It should be avoided using that value as an error code.
const NoCode ErrorCode = math.MaxUint16

// BadRequest can be used to signal user errors
const BadRequest ErrorCode = 400

// NotFound indicates unavailable resources
const NotFound ErrorCode = 404

type coder interface {
	Code() int
}

// GetCauseCode returns the attached code of the original causer of the error cascade.
func GetCauseCode(err error) int {

	cause := Cause(err)

	if errWithCode, ok := cause.(coder); ok {
		return errWithCode.Code()
	}

	return int(NoCode)
}

// Code traverses down the error cascade and return the first error with attached code it finds.
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
