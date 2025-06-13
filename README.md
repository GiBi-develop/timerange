# **TimeRange**
**Go-библиотека для работы с временными интервалами**  

Удобные операции для работы с промежутками времени:
- Проверка пересечений, объединение, поиск "окон".
- Форматирование в ISO 8601.
- Совместимо с `time.Time`.

---

## **Установка**
```bash
go get github.com/GiBi-develop/timerange@latest
```

---

## **Примеры использования**

### **1. Создание интервала**
```go
start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
end := start.Add(24 * time.Hour)
tr, err := timerange.New(start, end) // ошибка, если end < start
```

### **2. Проверка пересечения интервалов**
```go
other := timerange.TimeRange{
    Start: start.Add(12 * time.Hour),
    End:   end.Add(12 * time.Hour),
}
fmt.Println(tr.Overlaps(other)) // true
```

### **3. Объединение интервалов**
```go
merged, err := tr.Merge(other) // ошибка, если интервалы не пересекаются и не смежны
fmt.Println(merged.Duration()) // 36h0m0s
```

### **4. Поиск промежутка между интервалами**
```go
tr2, _ := timerange.New(
    time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
    time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC),
)
gap := tr.Gap(tr2) // промежуток между tr и tr2
fmt.Println(gap.ToISOString()) // "2023-01-02T00:00:00Z/2023-01-03T00:00:00Z"
```

### **5. Разделение интервала на части**
```go
parts := tr.SplitByDuration(6 * time.Hour) // []TimeRange
for _, part := range parts {
    fmt.Println(part.Start, "-", part.End)
}
```

---

## **Полный список методов**

| Метод | Описание | Пример |
|-------|----------|--------|
| `New(start, end time.Time)` | Создает интервал (проверяет валидность) | `tr, err := timerange.New(t1, t2)` |
| `Overlaps(other TimeRange)` | Проверяет пересечение интервалов | `tr.Overlaps(other)` |
| `Contains(t time.Time)` | Проверяет, содержит ли время | `tr.Contains(now)` |
| `Duration()` | Возвращает длительность | `tr.Duration()` |
| `SplitByDuration(d time.Duration)` | Делит интервал на части | `tr.SplitByDuration(time.Hour)` |
| `Merge(other TimeRange)` | Объединяет пересекающиеся интервалы | `merged, err := tr.Merge(other)` |
| `Gap(other TimeRange)` | Возвращает промежуток между интервалами | `gap := tr.Gap(other)` |
| `IsAdjacent(other TimeRange)` | Проверяет, смежны ли интервалы | `tr.IsAdjacent(other)` |
| `ToISOString()` | Форматирует в ISO 8601 | `tr.ToISOString()` |