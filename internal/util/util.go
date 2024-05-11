package util

import (
	"fmt"
	"io"
)

func HandleInputError(out io.Writer, line string) {
	fmt.Fprintf(out, "Invalid format found in line: '%s'", line)
}
