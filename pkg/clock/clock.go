package clock

import "time"

type Clock struct {
	NowFn func() time.Time
}

func (c Clock) Now() time.Time {
	if c.NowFn == nil {
		return time.Now()
	}
	return c.NowFn()
}
