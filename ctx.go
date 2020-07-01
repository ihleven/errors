package errors

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"
)

func Code(code int) *ctx {
	return &ctx{
		status: code,
	}
}

func Wrapnew(err error, msg string) error {
	c := &ctx{
		error:   err,
		message: msg,
	}
	file, function, line, ok := caller()
	if ok {
		c.file = file
		c.function = function
		c.line = line
	}
	return c
}

type ctx struct {
	error
	status  int
	message string
	file    string
	line    int
	// pack     string
	function string
	*stack
}

func (c *ctx) Error() string { return c.message }
func (c *ctx) Cause() error  { return c.error }
func (c *ctx) Unwrap() error { return c.error }
func (c *ctx) Dump() string {
	var str string
	newline := func() {
		if str != "" && !strings.HasSuffix(str, "\n") {
			str += "\n"
		}
	}
	for curr, ok := c, true; ok; curr, ok = curr.error.(*ctx) {
		str += curr.message
		if curr.file != "" {
			newline()
			if curr.function == "" {
				str += fmt.Sprintf(" --- at %v:%v ---", curr.file, curr.line)
			} else {
				str += fmt.Sprintf(" --- at %v:%v (%v) ---", curr.file, curr.line, curr.function)
			}
		}
		if curr.error != nil {
			newline()
			if cause, ok := curr.error.(*ctx); !ok {
				str += "Caused by: "
				str += curr.error.Error()
			} else if cause.message != "" {
				str += "Caused by: "
			}
		}
	}
	return str
}

func (c *ctx) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", c.message)
			// c.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, c.Dump())
	case 'q':
		fmt.Fprintf(s, "%q", c.Error())
	}
}

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func cleanpath(path string) string {

	rel, err := filepath.Rel(basepath, path)
	// filepath.Rel can traverse parent directories, don't want those
	if err == nil { //&& !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return rel
	}
	return ""
}
func caller() (file string, function string, line int, ok bool) {

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return
	}

	file = cleanpath(file)

	f := runtime.FuncForPC(pc)
	if f == nil {
		return
	}
	function = shortFuncName(f)

	return
}

/* "FuncName" or "Receiver.MethodName" */
func shortFuncName(f *runtime.Func) string {
	// f.Name() is like one of these:
	// - "github.com/palantir/shield/package.FuncName"
	// - "github.com/palantir/shield/package.Receiver.MethodName"
	// - "github.com/palantir/shield/package.(*PtrReceiver).MethodName"
	longName := f.Name()
	// return fmt.Sprintf("%s() %s", longName, os.Getenv("GOMOD"))
	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	// withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPath
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}
