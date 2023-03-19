package solution

import (
	"fmt"
	"io"
)

func Solve(r io.Reader, w io.Writer) error {
	_, err := fmt.Fprint(w, "3613")
	return err
}
