package sub1

import (
	"errdemo/sub1/sub2"

	"github.com/pkg/errors"
)

// Diff 比较差异
func Diff(foo, bar int) error {
	if foo < 0 {
		return errors.New("diff error")
	}

	if err := sub2.Diff(foo, bar); err != nil {
		return err
	}

	return nil
}
