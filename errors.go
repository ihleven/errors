package errors

import (
	"fmt"
	"io"
)

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
// In case a format string is given, New formats
// according to a format specifier and returns the string as a value that satisfies error.
func New(format string, args ...interface{}) error {

	return &withStack{
		&fundamental{
			msg:  fmt.Sprintf(format, args...),
			code: NoCode,
		},
		callers(),
	}
}

func NewWithCode(code ErrorCode, format string, args ...interface{}) error {

	return &withStack{
		&fundamental{
			msg:  fmt.Sprintf(format, args...),
			code: code,
		},
		callers(),
	}
}

// fundamental is an error that has a message and a stack, but no caller.
type fundamental struct {
	msg string
	// *stack
	code ErrorCode
}

func (f *fundamental) Error() string { return f.msg }

func (f *fundamental) Code() int { return int(f.code) }

func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		// if s.Flag('+') {
		// 	io.WriteString(s, f.msg)
		// 	// f.stack.Format(s, verb)
		// 	return
		// }
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
	}
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Cause() error { return w.error }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withStack) Unwrap() error { return w.error }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') || s.Flag('#') {
			fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

// Wrap returns an error annotating err with context and a stacktrace if
// * err is fundamental, withCode or unknown error
// if err is withStack or withContext, a stacktrace is already contained.

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	switch err.(type) {
	case *withStack, *withMessage:
	// nothing to do here
	default:
		// no stack as of yet, adding one
		err = &withStack{
			err,
			callers(),
		}
	}

	wrapped := &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	file, function, line, ok := caller()
	if ok {
		// fmt.Printf("wrapping:  %s:%d\n", file, line)
		wrapped.file = file
		wrapped.function = function
		wrapped.line = line
	}
	return wrapped
}

// Wrapf is an alias for Wrap to be compatible with pkg/errors.
func Wrapf(err error, format string, args ...interface{}) error {
	return Wrap(err, format, args...)
}

// WithMessage is an alias for Wrap to be compatible with pkg/errors.
func WithMessage(err error, format string, args ...interface{}) error {
	return Wrap(err, format, args...)
}

// WithMessagef is an alias for Wrap to be compatible with pkg/errors.
func WithMessagef(err error, format string, args ...interface{}) error {
	return Wrap(err, format, args...)
}

type withMessage struct {
	cause    error
	msg      string
	function string // function initiating withMessage
	file     string // file initiating withMessage
	line     int    // line initiating withMessage
}

func (w *withMessage) Error() string {
	if w.msg != "" {
		return w.msg + ": " + w.cause.Error()
	}
	return w.cause.Error()
}
func (w *withMessage) Cause() error { return w.cause }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withMessage) Unwrap() error { return w.cause }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':

		if s.Flag('+') {
			// fmt.Fprintf(s, "%+v\n", w.Cause())
			// io.WriteString(s, w.msg)
			// return
			io.WriteString(s, w.msg)
			fmt.Fprintf(s, "\n\t--- at %s:%d (%s)", w.file, w.line, w.function)
			if e, ok := w.Cause().(*withMessage); !ok || (ok && e.msg != "") {
				fmt.Fprintf(s, "\nCaused by: ")
			}
			fmt.Fprintf(s, "%+v\n", w.Cause())

			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
