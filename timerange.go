package timerange

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"
)

var (
	ErrInvalidRange    = errors.New("end time must be after start time")
	ErrNoOverlap       = errors.New("time ranges do not overlap")
	ErrNoIntersection  = errors.New("no intersection found")
	ErrInvalidArgument = errors.New("invalid argument")
)

type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// --- Core Functions ---

func New(start, end time.Time) (TimeRange, error) {
	if end.Before(start) {
		return TimeRange{}, ErrInvalidRange
	}
	return TimeRange{Start: start, End: end}, nil
}

// --- Basic Operations ---

func (tr TimeRange) Overlaps(other TimeRange) bool {
	return tr.Start.Before(other.End) && other.Start.Before(tr.End)
}

func (tr TimeRange) Contains(t time.Time) bool {
	return !t.Before(tr.Start) && !t.After(tr.End)
}

func (tr TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

// --- Set Operations ---

func Union(ranges []TimeRange) ([]TimeRange, error) {
	if len(ranges) == 0 {
		return nil, ErrInvalidArgument
	}

	sorted := make([]TimeRange, len(ranges))
	copy(sorted, ranges)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Start.Before(sorted[j].Start)
	})

	merged := []TimeRange{sorted[0]}
	for _, curr := range sorted[1:] {
		last := &merged[len(merged)-1]
		if curr.Start.Before(last.End) || curr.Start.Equal(last.End) {
			if curr.End.After(last.End) {
				last.End = curr.End
			}
		} else {
			merged = append(merged, curr)
		}
	}
	return merged, nil
}

func Intersection(ranges []TimeRange) (TimeRange, error) {
	if len(ranges) == 0 {
		return TimeRange{}, ErrInvalidArgument
	}

	maxStart := ranges[0].Start
	minEnd := ranges[0].End

	for _, tr := range ranges[1:] {
		if tr.Start.After(maxStart) {
			maxStart = tr.Start
		}
		if tr.End.Before(minEnd) {
			minEnd = tr.End
		}
	}

	if maxStart.After(minEnd) {
		return TimeRange{}, ErrNoIntersection
	}
	return TimeRange{Start: maxStart, End: minEnd}, nil
}

// --- Range Manipulation ---

func (tr TimeRange) SplitByDuration(d time.Duration) []TimeRange {
	if d <= 0 {
		return []TimeRange{tr}
	}

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

func (tr TimeRange) Merge(other TimeRange) (TimeRange, error) {
	if !tr.Overlaps(other) && !tr.IsAdjacent(other) {
		return TimeRange{}, ErrNoOverlap
	}
	return TimeRange{
		Start: minTime(tr.Start, other.Start),
		End:   maxTime(tr.End, other.End),
	}, nil
}

func (tr TimeRange) Subtract(other TimeRange) []TimeRange {
	if !tr.Overlaps(other) {
		return []TimeRange{tr}
	}

	var result []TimeRange
	if tr.Start.Before(other.Start) {
		result = append(result, TimeRange{
			Start: tr.Start,
			End:   other.Start,
		})
	}

	if other.End.Before(tr.End) {
		result = append(result, TimeRange{
			Start: other.End,
			End:   tr.End,
		})
	}
	return result
}

func (tr TimeRange) Gap(other TimeRange) TimeRange {
	if tr.Overlaps(other) || tr.IsAdjacent(other) {
		return TimeRange{}
	}
	if tr.End.Before(other.Start) {
		return TimeRange{Start: tr.End, End: other.Start}
	}
	return TimeRange{Start: other.End, End: tr.Start}
}

// --- Utility Functions ---

func (tr TimeRange) IsZero() bool {
	return tr.Start.IsZero() && tr.End.IsZero()
}

func (tr TimeRange) Equal(other TimeRange) bool {
	return tr.Start.Equal(other.Start) && tr.End.Equal(other.End)
}

func (tr TimeRange) IsAdjacent(other TimeRange) bool {
	return tr.End.Equal(other.Start) || other.End.Equal(tr.Start)
}

func (tr TimeRange) Clamp(t time.Time) time.Time {
	if t.Before(tr.Start) {
		return tr.Start
	}
	if t.After(tr.End) {
		return tr.End
	}
	return t
}

// --- Formatting ---

func (tr TimeRange) ToISOString() string {
	return fmt.Sprintf("%s/%s",
		tr.Start.Format(time.RFC3339),
		tr.End.Format(time.RFC3339),
	)
}

func (tr TimeRange) ToHumanString(layout string) string {
	if layout == "" {
		layout = time.RFC1123
	}
	return fmt.Sprintf("%s - %s",
		tr.Start.Format(layout),
		tr.End.Format(layout),
	)
}

func (tr TimeRange) ToSlugString() string {
	return fmt.Sprintf("%s-%s",
		tr.Start.Format("20060102"),
		tr.End.Format("20060102"),
	)
}

// --- JSON Support ---

func (tr TimeRange) MarshalJSON() ([]byte, error) {
	type Alias TimeRange
	return json.Marshal(&struct {
		Alias
		ISOString string `json:"iso"`
	}{
		Alias:     Alias(tr),
		ISOString: tr.ToISOString(),
	})
}

func (tr *TimeRange) UnmarshalJSON(data []byte) error {
	type Alias TimeRange
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(tr),
	}
	return json.Unmarshal(data, aux)
}

// --- Helper Functions ---

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
