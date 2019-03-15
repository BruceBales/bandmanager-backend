package structs

import "time"

type show struct {
	ID    int
	bands []Band
	date  time.Time
}
