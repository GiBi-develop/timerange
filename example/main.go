package main

import (
	"fmt"
	"github.com/GiBi-develop/timerange"
	"time"
)

func main() {
	start := time.Now()
	end := start.Add(2 * time.Hour)
	tr, err := timerange.New(start, end)
	if err != nil {
		panic(err)
	}

	// Проверка пересечения
	other := timerange.TimeRange{
		Start: start.Add(1 * time.Hour),
		End:   start.Add(3 * time.Hour),
	}
	fmt.Println("Overlaps:", tr.Overlaps(other)) // true

	// Разделение по интервалам
	ranges := tr.SplitByDuration(30 * time.Minute)
	for _, r := range ranges {
		fmt.Printf("%v - %v\n", r.Start, r.End)
	}
}
