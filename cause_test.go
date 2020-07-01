package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ihleven/errors"
)

type customError string

func (e customError) Error() string { return string(e) }

func TestCause(t *testing.T) {
	for _, test := range []struct {
		err   error
		cause error
	}{
		{
			err:   nil,
			cause: nil,
		},
		{
			err:   errors.New("msg"),
			cause: errors.New("msg"),
		},
		{
			err:   errors.New("msg"),
			cause: errors.New("msg"),
		},
		{
			err:   errors.Wrap(errors.New("msg1"), "msg2"),
			cause: errors.New("msg1"),
		},
		{
			err:   customError("msg"),
			cause: customError("msg"),
		},
		{
			err:   errors.Wrap(customError("msg1"), "msg2"),
			cause: customError("msg1"),
		},
	} {
		assert.Equal(t, test.cause, errors.Cause(test.err))
	}
}
