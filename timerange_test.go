package timerange

import (
	"testing"
	"time"
)

func TestMerge(t *testing.T) {
	tr1, _ := New(
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	)
	tr2, _ := New(
		time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
	)

	merged, err := tr1.Merge(tr2)
	if err != nil {
		t.Fatal(err)
	}
	expectedStart := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)

	if !merged.Start.Equal(expectedStart) || !merged.End.Equal(expectedEnd) {
		t.Errorf("Merge() = %v, want %v to %v", merged, expectedStart, expectedEnd)
	}
}

func TestGap(t *testing.T) {
	tr1, _ := New(
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	)
	tr2, _ := New(
		time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC),
	)

	gap := tr1.Gap(tr2)
	expectedStart := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)

	if !gap.Start.Equal(expectedStart) || !gap.End.Equal(expectedEnd) {
		t.Errorf("Gap() = %v, want %v to %v", gap, expectedStart, expectedEnd)
	}
}

func TestToISOString(t *testing.T) {
	tr, _ := New(
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	)
	expected := "2023-01-01T00:00:00Z/2023-01-02T00:00:00Z"
	if tr.ToISOString() != expected {
		t.Errorf("ToISOString() = %v, want %v", tr.ToISOString(), expected)
	}
}
