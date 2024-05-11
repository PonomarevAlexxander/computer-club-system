package stats

import (
	"math"
	"time"
)

type TableStats struct {
	Income    uint
	TimeTaken time.Time
}

func (t *TableStats) AddIncomeUpperBound(duration time.Duration, hourPrice uint) {
	t.Income += hourPrice * uint(math.Ceil(duration.Hours()))
}

func (t *TableStats) AddIncomeLowerBound(duration time.Duration, hourPrice uint) {
	t.Income += hourPrice * uint(math.Floor(duration.Hours()))
}

func (t *TableStats) AddTime(duration time.Duration) {
	t.TimeTaken = t.TimeTaken.Add(duration)
}
