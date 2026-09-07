package structs

import "time"

type LinkStat struct {
	LinkID      int64
	StatDate    time.Time
	ClicksCount int64
}
