package timerange

import (
	"errors"
	"time"
)

var ErrInvalidRange = errors.New("end time must be after start time")

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func New(start, end time.Time) (TimeRange, error) {
	if end.Before(start) {
		return TimeRange{}, ErrInvalidRange
	}
	return TimeRange{Start: start, End: end}, nil
}

func (tr TimeRange) Overlaps(other TimeRange) bool {
	return tr.Start.Before(other.End) && other.Start.Before(tr.End)
}

func (tr TimeRange) Contains(t time.Time) bool {
	return !t.Before(tr.Start) && !t.After(tr.End)
}

func (tr TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

func (tr TimeRange) SplitByDuration(d time.Duration) []TimeRange {
	var ranges []TimeRange
	current := tr.Start

	for current.Before(tr.End) {
		next := current.Add(d)
		if next.After(tr.End) {
			next = tr.End
		}
		ranges = append(ranges, TimeRange{Start: current, End: next})
		current = next
	}

	return ranges
}
