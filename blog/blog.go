package blog

import (
	"fmt"
	"log"
	"os"

	"github.com/rotisserie/eris"
)

// Verbosity is a relative indicator of how much detail we should provide to
// others when accepting log messages.
var Verbosity = 0

func Info(spec string, args ...any) {
	if Verbosity > 0 {
		log.Printf(spec, args...)
	}
}

func Fail(err error) {
	fmt.Fprintln(os.Stderr, eris.ToString(err, Verbosity > 0))
	os.Exit(1)
}

func Warn(err error) {
	fmt.Fprintln(os.Stderr, eris.ToString(err, Verbosity > 0))
}
