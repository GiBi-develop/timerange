package timerange

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidRange = errors.New("end time must be after start time")
	ErrNoOverlap    = errors.New("time ranges do not overlap")
)

type TimeRange struct {
	Start time.Time
	End   time.Time
}

// New создаёт TimeRange. Возвращает ошибку, если end < start.
func New(start, end time.Time) (TimeRange, error) {
	if end.Before(start) {
		return TimeRange{}, ErrInvalidRange
	}
	return TimeRange{Start: start, End: end}, nil
}

// Overlaps проверяет, пересекаются ли интервалы.
func (tr TimeRange) Overlaps(other TimeRange) bool {
	return tr.Start.Before(other.End) && other.Start.Before(tr.End)
}

// Contains проверяет, содержит ли интервал указанное время.
func (tr TimeRange) Contains(t time.Time) bool {
	return !t.Before(tr.Start) && !t.After(tr.End)
}

// Duration возвращает длительность интервала.
func (tr TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

// SplitByDuration делит интервал на подынтервалы длительности d.
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

// Merge объединяет два пересекающихся интервала.
func (tr TimeRange) Merge(other TimeRange) (TimeRange, error) {
	if !tr.Overlaps(other) && !tr.IsAdjacent(other) {
		return TimeRange{}, ErrNoOverlap
	}
	start := minTime(tr.Start, other.Start)
	end := maxTime(tr.End, other.End)
	return TimeRange{Start: start, End: end}, nil
}

// Gap возвращает промежуток между непересекающимися интервалами.
func (tr TimeRange) Gap(other TimeRange) TimeRange {
	if tr.Overlaps(other) || tr.IsAdjacent(other) {
		return TimeRange{}
	}
	if tr.End.Before(other.Start) {
		return TimeRange{Start: tr.End, End: other.Start}
	}
	return TimeRange{Start: other.End, End: tr.Start}
}

// IsAdjacent проверяет, являются ли интервалы смежными.
func (tr TimeRange) IsAdjacent(other TimeRange) bool {
	return tr.End.Equal(other.Start) || other.End.Equal(tr.Start)
}

// ToISOString возвращает строку в формате ISO 8601 (например, "2006-01-02T15:04:05Z/2006-01-03T15:04:05Z").
func (tr TimeRange) ToISOString() string {
	return fmt.Sprintf("%s/%s", tr.Start.Format(time.RFC3339), tr.End.Format(time.RFC3339))
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
