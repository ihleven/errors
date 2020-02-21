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
	return errors.Wrap(e, "")
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
	return errors.New("This is a new error")
}

func Example() {
	e := wrappedError()
	fmt.Printf("%+#v\n", e)
	// Output:
	// 	Failed to wrap accordingly

}
