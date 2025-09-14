package xrrtest

import (
	"time"

	"github.com/ctx42/xrr/pkg/xrr"
)

// TstError returns a test error with meta data.
func TstError() error {
	tim := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

	e := xrr.New("msg", "EC", xrr.Meta().Int("A", 1).Str("str", "abc").Option())
	e = xrr.Wrap(e, xrr.Meta().Int("A", 2).Int("int", 2).Option())
	e = xrr.Wrap(e, xrr.Meta().Int("A", 3).Int64("int64", 3).Option())
	e = xrr.Wrap(e, xrr.Meta().Int("A", 4).Float64("float64", 4).Option())
	e = xrr.Wrap(e, xrr.Meta().Int("A", 5).Bool("bool", true).Option())
	e = xrr.Wrap(e, xrr.Meta().Int("A", 6).Time("tim", tim).Option())

	return e
}
