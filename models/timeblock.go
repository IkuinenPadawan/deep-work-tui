package models

import "time"

type Timeblock struct {
	Task      string
	Starttime time.Time
	Endtime   time.Time
}
