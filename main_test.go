package errors_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ihleven/errors"
)

func TestCaller(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	fmt.Println(b, basepath)
}
func TestMain(t *testing.T) {
	fmt.Println("")

	e := wrapper()
	fmt.Printf("%+v\n", e)
}

// func main() {
// 	stacktrace.Propagate(wrapper(), "und das der Kontext")
// 	// fmt.Printf("%+v\n", err)
// }
func wrapper() error {
	e := caused()
	return errors.Wrap(e, "Failed to wrap accordingly")
	// return errors.New(0, "kjhlkhlkjhlj")
}
func caused() error {
	e := ohne()
	// e := errors.New(0, "kjhlkhlkjhlj")
	return errors.Wrap(e) //, "Couldn't cause as expected!")
}
func ohne() error {
	e := grund()
	return errors.Wrap(e, "could not find filename")

}

func grund() error {
	return errors.New("Hallo %s", "Welt")
	// _, err := os.Open("kjhklh")
	// // fmt.Printf("%+v\n", e)
	// return errors.Wrap(err, "Wrapping os open error")
}

func Nuller() error {
	var blah interface{ Get() (int, error) }
	if _, err := blah.Get(); err != nil {
		return err
	}
	return nil
}
