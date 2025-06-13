package timerange

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	start := time.Now()
	end := start.Add(1 * time.Hour)

	_, err := New(start, end)
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}

	_, err = New(end, start) // Неправильный порядок
	if err == nil {
		t.Fatal("New() with invalid range, want error")
	}
}

func TestOverlaps(t *testing.T) {
	tr1, _ := New(
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	)
	tr2, _ := New(
		time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
	)

	if !tr1.Overlaps(tr2) {
		t.Error("Overlaps() = false, want true")
	}
}
