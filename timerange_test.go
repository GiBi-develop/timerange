package timerange

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("valid range", func(t *testing.T) {
		start := time.Now()
		end := start.Add(time.Hour)
		_, err := New(start, end)
		if err != nil {
			t.Errorf("New() error = %v, want nil", err)
		}
	})

	t.Run("invalid range", func(t *testing.T) {
		end := time.Now()
		start := end.Add(time.Hour)
		_, err := New(start, end)
		if err != ErrInvalidRange {
			t.Errorf("New() error = %v, want ErrInvalidRange", err)
		}
	})

	t.Run("zero-length range", func(t *testing.T) {
		now := time.Now()
		tr, err := New(now, now)
		if err != nil {
			t.Errorf("New() error = %v, want nil", err)
		}
		if !tr.IsZero() {
			t.Error("IsZero() = false, want true for zero-length range")
		}
	})
}

func TestZeroValue(t *testing.T) {
	var tr TimeRange
	if !tr.IsZero() {
		t.Error("IsZero() should return true for zero value TimeRange")
	}
}

func TestOverlaps(t *testing.T) {
	base := TimeRange{
		Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name   string
		other  TimeRange
		expect bool
	}{
		{
			name: "complete overlap",
			other: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			expect: true,
		},
		{
			name: "partial overlap",
			other: TimeRange{
				Start: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			expect: true,
		},
		{
			name: "no overlap",
			other: TimeRange{
				Start: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC),
			},
			expect: false,
		},
		{
			name: "adjacent",
			other: TimeRange{
				Start: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := base.Overlaps(tt.other)
			if result != tt.expect {
				t.Errorf("Overlaps() = %v, want %v", result, tt.expect)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tr := TimeRange{
		Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name   string
		t      time.Time
		expect bool
	}{
		{"before start", time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC), false},
		{"at start", tr.Start, true},
		{"middle", time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), true},
		{"at end", tr.End, true},
		{"after end", time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tr.Contains(tt.t)
			if result != tt.expect {
				t.Errorf("Contains() = %v, want %v", result, tt.expect)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	start := time.Now()
	end := start.Add(2 * time.Hour)
	tr := TimeRange{Start: start, End: end}

	if tr.Duration() != 2*time.Hour {
		t.Errorf("Duration() = %v, want %v", tr.Duration(), 2*time.Hour)
	}
}

func TestSplitByDuration(t *testing.T) {
	tr := TimeRange{
		Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	t.Run("exact division", func(t *testing.T) {
		parts := tr.SplitByDuration(12 * time.Hour)
		if len(parts) != 2 {
			t.Fatalf("SplitByDuration() len = %d, want 2", len(parts))
		}
	})

	t.Run("uneven division", func(t *testing.T) {
		parts := tr.SplitByDuration(18 * time.Hour)
		if len(parts) != 2 {
			t.Fatalf("SplitByDuration() len = %d, want 2", len(parts))
		}
	})

	t.Run("larger than range", func(t *testing.T) {
		parts := tr.SplitByDuration(48 * time.Hour)
		if len(parts) != 1 {
			t.Fatalf("SplitByDuration() len = %d, want 1", len(parts))
		}
	})
}

func TestMerge(t *testing.T) {
	t.Run("overlapping", func(t *testing.T) {
		tr1 := TimeRange{
			Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		}
		tr2 := TimeRange{
			Start: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
		}

		merged, err := tr1.Merge(tr2)
		if err != nil {
			t.Fatal(err)
		}
		expected := TimeRange{
			Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
		}
		if !merged.Equal(expected) {
			t.Errorf("Merge() = %v, want %v", merged, expected)
		}
	})

	t.Run("non-overlapping", func(t *testing.T) {
		tr1 := TimeRange{
			Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		}
		tr2 := TimeRange{
			Start: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC),
		}

		_, err := tr1.Merge(tr2)
		if err != ErrNoOverlap {
			t.Errorf("Merge() error = %v, want ErrNoOverlap", err)
		}
	})
}

func TestUnion(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		_, err := Union([]TimeRange{})
		if err != ErrInvalidArgument {
			t.Errorf("Union() error = %v, want ErrInvalidArgument", err)
		}
	})

	t.Run("multiple ranges", func(t *testing.T) {
		ranges := []TimeRange{
			{Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
			{Start: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), End: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)},
			{Start: time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC), End: time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC)},
		}

		result, err := Union(ranges)
		if err != nil {
			t.Fatal(err)
		}
		if len(result) != 2 {
			t.Fatalf("Union() len = %d, want 2", len(result))
		}
	})
}

func TestIntersection(t *testing.T) {
	t.Run("no intersection", func(t *testing.T) {
		ranges := []TimeRange{
			{Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
			{Start: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), End: time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC)},
		}

		_, err := Intersection(ranges)
		if err != ErrNoIntersection {
			t.Errorf("Intersection() error = %v, want ErrNoIntersection", err)
		}
	})

	t.Run("valid intersection", func(t *testing.T) {
		ranges := []TimeRange{
			{Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)},
			{Start: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), End: time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC)},
		}

		result, err := Intersection(ranges)
		if err != nil {
			t.Fatal(err)
		}
		expected := TimeRange{
			Start: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
		}
		if !result.Equal(expected) {
			t.Errorf("Intersection() = %v, want %v", result, expected)
		}
	})
}

func TestJSON(t *testing.T) {
	tr := TimeRange{
		Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(tr)
	if err != nil {
		t.Fatal(err)
	}

	var unmarshaled TimeRange
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatal(err)
	}

	if !tr.Equal(unmarshaled) {
		t.Errorf("JSON roundtrip failed: got %v, want %v", unmarshaled, tr)
	}
}

func TestGap(t *testing.T) {
	tr := TimeRange{
		Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name   string
		other  TimeRange
		expect TimeRange
	}{
		{
			name: "before with gap",
			other: TimeRange{
				Start: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC),
			},
			expect: TimeRange{
				Start: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "after with gap",
			other: TimeRange{
				Start: time.Date(2022, 12, 30, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC),
			},
			expect: TimeRange{
				Start: time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "overlapping",
			other: TimeRange{
				Start: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			expect: TimeRange{},
		},
		{
			name: "adjacent",
			other: TimeRange{
				Start: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			expect: TimeRange{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tr.Gap(tt.other)
			if !result.Equal(tt.expect) {
				t.Errorf("Gap() = %v, want %v", result, tt.expect)
			}
		})
	}
}

func TestMergeOverlapping(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		result, err := MergeOverlapping([]TimeRange{})
		if err != nil {
			t.Fatal(err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty result, got %v", result)
		}
	})

	t.Run("non-overlapping", func(t *testing.T) {
		input := []TimeRange{
			{Start: parseTime("2023-01-01"), End: parseTime("2023-01-02")},
			{Start: parseTime("2023-01-03"), End: parseTime("2023-01-04")},
		}
		result, err := MergeOverlapping(input)
		if err != nil {
			t.Fatal(err)
		}
		if len(result) != 2 {
			t.Errorf("Expected 2 ranges, got %d", len(result))
		}
	})

	t.Run("overlapping", func(t *testing.T) {
		input := []TimeRange{
			{Start: parseTime("2023-01-01"), End: parseTime("2023-01-03")},
			{Start: parseTime("2023-01-02"), End: parseTime("2023-01-04")},
		}
		expected := TimeRange{
			Start: parseTime("2023-01-01"),
			End:   parseTime("2023-01-04"),
		}
		result, err := MergeOverlapping(input)
		if err != nil {
			t.Fatal(err)
		}
		if len(result) != 1 || !result[0].Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestFindGaps(t *testing.T) {
	bounds := TimeRange{
		Start: parseTime("2023-01-01"),
		End:   parseTime("2023-01-10"),
	}

	t.Run("no occupied", func(t *testing.T) {
		gaps, err := FindGaps([]TimeRange{}, bounds)
		if err != nil {
			t.Fatal(err)
		}
		if len(gaps) != 1 || !gaps[0].Equal(bounds) {
			t.Errorf("Expected full bounds as gap, got %v", gaps)
		}
	})

	t.Run("with gaps", func(t *testing.T) {
		occupied := []TimeRange{
			{Start: parseTime("2023-01-02"), End: parseTime("2023-01-03")},
			{Start: parseTime("2023-01-05"), End: parseTime("2023-01-07")},
		}
		expected := []TimeRange{
			{Start: parseTime("2023-01-01"), End: parseTime("2023-01-02")},
			{Start: parseTime("2023-01-03"), End: parseTime("2023-01-05")},
			{Start: parseTime("2023-01-07"), End: parseTime("2023-01-10")},
		}
		gaps, err := FindGaps(occupied, bounds)
		if err != nil {
			t.Fatal(err)
		}
		if !compareRanges(gaps, expected) {
			t.Errorf("Expected %v, got %v", expected, gaps)
		}
	})
}

func TestIsAdjacent(t *testing.T) {
	tr := TimeRange{
		Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name   string
		other  TimeRange
		expect bool
	}{
		{
			name: "left adjacent",
			other: TimeRange{
				Start: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			expect: true,
		},
		{
			name: "right adjacent",
			other: TimeRange{
				Start: time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expect: true,
		},
		{
			name: "not adjacent",
			other: TimeRange{
				Start: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC),
			},
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tr.IsAdjacent(tt.other)
			if result != tt.expect {
				t.Errorf("IsAdjacent() = %v, want %v", result, tt.expect)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tr := TimeRange{
		Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name   string
		other  TimeRange
		expect []TimeRange
	}{
		{
			name: "no overlap",
			other: TimeRange{
				Start: time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC),
			},
			expect: []TimeRange{tr},
		},
		{
			name: "full overlap",
			other: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			expect: []TimeRange{},
		},
		{
			name: "left overlap",
			other: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			expect: []TimeRange{
				{
					Start: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "middle overlap",
			other: TimeRange{
				Start: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			expect: []TimeRange{
				{
					Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				},
				{
					Start: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tr.Subtract(tt.other)
			if len(result) != len(tt.expect) {
				t.Fatalf("Subtract() len = %d, want %d", len(result), len(tt.expect))
			}
			for i := range result {
				if !result[i].Equal(tt.expect[i]) {
					t.Errorf("Subtract()[%d] = %v, want %v", i, result[i], tt.expect[i])
				}
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		name   string
		tr     TimeRange
		expect bool
	}{
		{
			name:   "zero range",
			tr:     TimeRange{},
			expect: true,
		},
		{
			name: "non-zero range",
			tr: TimeRange{
				Start: time.Now(),
				End:   time.Now().Add(time.Hour),
			},
			expect: false,
		},
		{
			name: "zero start only",
			tr: TimeRange{
				Start: time.Time{},
				End:   time.Now(),
			},
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tr.IsZero()
			if result != tt.expect {
				t.Errorf("IsZero() = %v, want %v", result, tt.expect)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	now := time.Now()
	tr := TimeRange{
		Start: now,
		End:   now.Add(time.Hour),
	}

	tests := []struct {
		name   string
		other  TimeRange
		expect bool
	}{
		{
			name: "equal",
			other: TimeRange{
				Start: now,
				End:   now.Add(time.Hour),
			},
			expect: true,
		},
		{
			name: "different start",
			other: TimeRange{
				Start: now.Add(-time.Hour),
				End:   now.Add(time.Hour),
			},
			expect: false,
		},
		{
			name: "different end",
			other: TimeRange{
				Start: now,
				End:   now.Add(2 * time.Hour),
			},
			expect: false,
		},
		{
			name:   "zero value",
			other:  TimeRange{},
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tr.Equal(tt.other)
			if result != tt.expect {
				t.Errorf("Equal() = %v, want %v", result, tt.expect)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tr := TimeRange{
		Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name   string
		t      time.Time
		expect time.Time
	}{
		{
			name:   "before range",
			t:      time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC),
			expect: tr.Start,
		},
		{
			name:   "in range",
			t:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			expect: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:   "after range",
			t:      time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			expect: tr.End,
		},
		{
			name:   "at start",
			t:      tr.Start,
			expect: tr.Start,
		},
		{
			name:   "at end",
			t:      tr.End,
			expect: tr.End,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tr.Clamp(tt.t)
			if !result.Equal(tt.expect) {
				t.Errorf("Clamp() = %v, want %v", result, tt.expect)
			}
		})
	}
}

func TestToStringMethods(t *testing.T) {
	tr := TimeRange{
		Start: time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
		End:   time.Date(2023, 1, 2, 13, 45, 0, 0, time.UTC),
	}

	t.Run("ToISOString", func(t *testing.T) {
		result := tr.ToISOString()
		expected := "2023-01-01T12:30:00Z/2023-01-02T13:45:00Z"
		if result != expected {
			t.Errorf("ToISOString() = %v, want %v", result, expected)
		}
	})

	t.Run("ToHumanString default layout", func(t *testing.T) {
		result := tr.ToHumanString("")
		expected := "Sun, 01 Jan 2023 12:30:00 UTC - Mon, 02 Jan 2023 13:45:00 UTC"
		if result != expected {
			t.Errorf("ToHumanString() = %v, want %v", result, expected)
		}
	})

	t.Run("ToHumanString", func(t *testing.T) {
		t.Run("default layout", func(t *testing.T) {
			result := tr.ToHumanString("")
			expected := "Sun, 01 Jan 2023 12:30:00 UTC - Mon, 02 Jan 2023 13:45:00 UTC"
			if result != expected {
				t.Errorf("ToHumanString() = %v, want %v", result, expected)
			}
		})

		t.Run("custom layout", func(t *testing.T) {
			result := tr.ToHumanString("2006-01-02 15:04")
			expected := "2023-01-01 12:30 - 2023-01-02 13:45"
			if result != expected {
				t.Errorf("ToHumanString() = %v, want %v", result, expected)
			}
		})
	})

	t.Run("ToSlugString", func(t *testing.T) {
		result := tr.ToSlugString()
		expected := "20230101-20230102"
		if result != expected {
			t.Errorf("ToSlugString() = %v, want %v", result, expected)
		}
	})
}

func TestJSONMarshaling(t *testing.T) {
	tr := TimeRange{
		Start: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		End:   time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
	}

	t.Run("MarshalJSON", func(t *testing.T) {
		data, err := json.Marshal(tr)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Unmarshal error = %v", err)
		}

		if parsed["start"] != "2023-01-01T12:00:00Z" {
			t.Errorf("MarshalJSON() start = %v, want %v", parsed["start"], "2023-01-01T12:00:00Z")
		}
		if parsed["iso"] != "2023-01-01T12:00:00Z/2023-01-02T12:00:00Z" {
			t.Errorf("MarshalJSON() iso = %v, want %v", parsed["iso"], "2023-01-01T12:00:00Z/2023-01-02T12:00:00Z")
		}
	})

	t.Run("Unmarshal empty object", func(t *testing.T) {
		var tr TimeRange
		err := json.Unmarshal([]byte(`{}`), &tr)
		if err == nil {
			t.Error("Expected error for empty object unmarshaling")
		}
	})

	t.Run("UnmarshalError", func(t *testing.T) {
		tests := []struct {
			name string
			data string
		}{
			{
				name: "invalid start",
				data: `{"start":"invalid","end":"2023-01-02T12:00:00Z"}`,
			},
			{
				name: "invalid end",
				data: `{"start":"2023-01-01T12:00:00Z","end":"invalid"}`,
			},
			{
				name: "empty object",
				data: `{}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var result TimeRange
				if err := json.Unmarshal([]byte(tt.data), &result); err == nil {
					t.Error("UnmarshalJSON() error = nil, want error")
				}
			})
		}
	})
}

func parseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func compareRanges(a, b []TimeRange) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}
