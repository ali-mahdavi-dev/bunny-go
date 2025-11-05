package errors

import "github.com/pkg/errors"

// Define alias for error wrapping utilities
var (
	WithStack = errors.WithStack
	Wrap      = errors.Wrap
	Wrapf     = errors.Wrapf
	Is        = errors.Is
	Errorf    = errors.Errorf
)

// As finds the first error in err's chain that matches Error
func As(err error) (Error, bool) {
	if err == nil {
		return nil, false
	}
	var appErr Error
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
