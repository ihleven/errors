package errors_test

import (
	"fmt"
	"os"

	"github.com/ihleven/errors"
)

func wrappedError() error {
	e := wrap()
	return errors.Wrap(e, "Failed to wrap accordingly")
}
func wrap() error {
	e := wrap2()
	return errors.Wrap(e, "Couldn't cause as expected!")
}
func wrap2() error {
	e := wrap3()
	return errors.Wrap(e, "asdf %d %d", 5, 6)
}

func wrap3() error {
	e := newerror()
	return errors.Wrap(e, "Could not find filename")
}

func open() error {
	_, err := os.Open("non-existing-filename")
	return errors.Wrap(err, "Could not open file")
}

func newerror() error {
	return errors.New("This is a new error %d", 6)
}

func Example() {
	e := wrappedError()
	fmt.Printf("%+#v\n", e)
	// Output:
	// Failed to wrap accordingly
	// 	--- at /Users/ih/src/errors/example_test.go:12 (wrappedError)
	// Caused by: Couldn't cause as expected!
	// 	--- at /Users/ih/src/errors/example_test.go:16 (wrap)
	// Caused by: asdf 5 6
	// 	--- at /Users/ih/src/errors/example_test.go:20 (wrap2)
	// Caused by: Could not find filename
	// 	--- at /Users/ih/src/errors/example_test.go:25 (wrap3)
	// Caused by: This is a new error 6
	// github.com/ihleven/errors_test.newerror
	// 	/Users/ih/src/errors/example_test.go:34
	// github.com/ihleven/errors_test.wrap3
	// 	/Users/ih/src/errors/example_test.go:24
	// github.com/ihleven/errors_test.wrap2
	// 	/Users/ih/src/errors/example_test.go:19
	// github.com/ihleven/errors_test.wrap
	// 	/Users/ih/src/errors/example_test.go:15
	// github.com/ihleven/errors_test.wrappedError
	// 	/Users/ih/src/errors/example_test.go:11
	// github.com/ihleven/errors_test.Example
	// 	/Users/ih/src/errors/example_test.go:38
	// testing.runExample
	// 	/usr/local/opt/go/libexec/src/testing/run_example.go:62
	// testing.runExamples
	// 	/usr/local/opt/go/libexec/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/opt/go/libexec/src/testing/testing.go:1200
	// main.main
	// 	_testmain.go:46
	// runtime.main
	// 	/usr/local/opt/go/libexec/src/runtime/proc.go:203
	// runtime.goexit
	// 	/usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1373


}
