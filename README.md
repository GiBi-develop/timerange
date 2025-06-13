# **TimeRange**
**Go-библиотека для работы с временными интервалами**  
[![Go Reference](https://pkg.go.dev/badge/github.com/GiBi-develop/timerange.svg)](https://pkg.go.dev/github.com/GiBi-develop/timerange)  
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Мощные операции для работы с временными промежутками:
- ✔️ Объединение, пересечение и вычитание интервалов
- ✔️ Поиск "окон" и смежных периодов
- ✔️ Поддержка JSON и различных форматов вывода
- ✔️ Полная совместимость с `time.Time`

---

## **Установка**
```bash
go get github.com/GiBi-develop/timerange@latest
```

---

## **Примеры использования**

### **1. Создание и базовые операции**
```go
// Создание интервала
tr, _ := timerange.New(
    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
    time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
)

// Проверка пересечения
other := timerange.TimeRange{
    Start: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
    End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
}
fmt.Println(tr.Overlaps(other)) // true
```

### **2. Работа с множеством интервалов**
```go
// Объединение пересекающихся интервалов
ranges := []timerange.TimeRange{tr1, tr2, tr3}
merged, _ := timerange.Union(ranges)

// Поиск общего времени
common, _ := timerange.Intersection(ranges)

// Вычитание интервалов
available := busyRange.Subtract(meetingRange)
```

### **3. Форматирование и сериализация**
```go
// Человекочитаемый формат
fmt.Println(tr.ToHumanString("2006-01-02")) // "2023-01-01 - 2023-01-02"

// JSON
jsonData, _ := json.Marshal(tr)
var restored timerange.TimeRange
json.Unmarshal(jsonData, &restored)
```

---

## **Полный справочник методов**

### **Создание и валидация**
| Метод | Описание | Пример |
|-------|----------|--------|
| `New(start, end time.Time)` | Создает новый интервал | `tr, err := timerange.New(start, end)` |
| `FromDuration(start time.Time, d time.Duration)` | Создает из начальной точки и длительности | `tr := timerange.FromDuration(now, 2*time.Hour)` |

### **Основные операции**
| Метод | Описание | Пример |
|-------|----------|--------|
| `Overlaps(other TimeRange)` | Проверяет пересечение | `if tr1.Overlaps(tr2)` |
| `Contains(t time.Time)` | Проверяет вхождение времени | `if tr.Contains(now)` |
| `Duration()` | Возвращает длительность | `dur := tr.Duration()` |
| `IsZero()` | Проверяет нулевой интервал | `if tr.IsZero()` |
| `Equal(other TimeRange)` | Сравнивает интервалы | `if tr1.Equal(tr2)` |

### **Операции с множествами**
| Метод | Описание | Пример |
|-------|----------|--------|
| `Union(ranges []TimeRange)` | Объединяет интервалы | `merged, _ := timerange.Union(ranges)` |
| `Intersection(ranges []TimeRange)` | Находит пересечение | `common, _ := timerange.Intersection(ranges)` |
| `Subtract(other TimeRange)` | Вычитает интервал | `remaining := tr.Subtract(other)` |
| `Gap(other TimeRange)` | Находит промежуток между интервалами | `gap := tr1.Gap(tr2)` |
| `MergeOverlapping(ranges []TimeRange)` | Объединяет пересекающиеся интервалы | `merged, _ := timerange.MergeOverlapping(ranges)` |
| `FindGaps(occupied []TimeRange, bounds TimeRange)` | Находит свободные промежутки | `gaps, _ := timerange.FindGaps(busy, dayBounds)` |


### **Форматирование**
| Метод | Описание | Пример |
|-------|----------|--------|
| `ToISOString()` | Формат ISO 8601 | `str := tr.ToISOString()` |
| `ToHumanString(layout)` | Читаемый формат | `str := tr.ToHumanString("Jan 2, 2006")` |
| `ToSlugString()` | Для URL и идентификаторов | `slug := tr.ToSlugString()` |

### **Утилиты**
| Метод | Описание | Пример |
|-------|----------|--------|
| `SplitByDuration(d)` | Делит на подынтервалы | `parts := tr.SplitByDuration(time.Hour)` |
| `Clamp(t time.Time)` | Ограничивает время интервалом | `safeTime := tr.Clamp(userTime)` |
| `IsAdjacent(other)` | Проверяет смежность | `if tr1.IsAdjacent(tr2)` |

---

## **Дополнительные примеры**

### **Поиск свободных окон**
```go
workDay, _ := timerange.New(
    time.Date(2023, 1, 1, 9, 0, 0, 0, time.UTC),
    time.Date(2023, 1, 1, 18, 0, 0, 0, time.UTC),
)

meetings := []timerange.TimeRange{meeting1, meeting2}
freeSlots := workDay.Subtract(timerange.Union(meetings))
```

### **Сериализация в JSON**
```go
type Event struct {
    Name      string          `json:"name"`
    TimeRange timerange.TimeRange `json:"time_range"`
}

event := Event{
    Name: "Meeting",
    TimeRange: tr,
}
jsonData, _ := json.Marshal(event)
```

### Объединение пересекающихся периодов
```go
ranges := []timerange.TimeRange{
    {Start: parse("2023-01-01"), End: parse("2023-01-03")},
    {Start: parse("2023-01-02"), End: parse("2023-01-04")},
}
merged, _ := timerange.MergeOverlapping(ranges)
// Результат: [{2023-01-01 00:00:00 +0000 UTC 2023-01-04 00:00:00 +0000 UTC}]
```
---

## **Лицензия**
MIT License. Полный текст доступен в файле [LICENSE](LICENSE).