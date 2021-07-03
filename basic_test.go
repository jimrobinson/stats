package stats

import (
	"fmt"
	"math/rand"
	"testing"
)

type statTest struct {
	Slots    SlotsFloat64
	Expected Stats
}

var statTests = []statTest{
	{SlotsFloat64{1, 1}, Stats{2, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1}},
	{SlotsFloat64{2, 1}, Stats{2, 1.5, 0.7071, 0.5, 1, 2, 1, 2, 2, 2, 2}},
	{SlotsFloat64{4, 3, 2, 1}, Stats{4, 2.5, 1.291, 1.6667, 1, 4, 2, 3, 4, 4, 4}},
	{SlotsFloat64{256, 1, 1024, 512}, Stats{4, 448.25, 436.8618, 190848.25, 1, 1024, 256, 512, 1024, 1024, 1024}},
	{SlotsFloat64{2, 1, 5, 4, 3, 6, 9, 8, 7, 10}, Stats{10, 5.5, 3.0277, 9.1667, 1, 10, 3, 6, 8, 10, 10}}}

func TestStatsF64(t *testing.T) {

	// test zero length array
	slots := make(SlotsFloat64, 0)
	stats, err := slots.Stats()
	if err != EMPTY {
		t.Errorf("expected error value of %v, got %v", EMPTY, err)
	}

	// test valid arrays
	for i := range statTests {
		slots := make(SlotsFloat64, len(statTests[i].Slots))
		copy(slots, statTests[i].Slots)

		stats, err = slots.Stats()
		if err != nil {
			t.Errorf("unexpected error value on statTests[%d].Slots.Stats(): %v", i, err)
		}
		err = equalStats(&statTests[i].Expected, stats)
		if err != nil {
			t.Error(err)
		}

		// same sort order?
		for j := range statTests[i].Slots {
			if slots[j] != statTests[i].Slots[j] {
				t.Error("slots have been reordered")
			}
		}
	}
}

func BenchmarkStats1K(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		slots := make(SlotsFloat64, 1<<10)
		for i := range slots {
			slots[i] = rand.Float64()
		}
		b.StartTimer()
		slots.Stats()
		b.StopTimer()
	}
}

func BenchmarkStats64K(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		slots := make(SlotsFloat64, 1<<16)
		for i := range slots {
			slots[i] = rand.Float64()
		}
		b.StartTimer()
		slots.Stats()
		b.StopTimer()
	}
}

// equalStats examines Stats a and b, and returns an error if any of
// the members do not match when rounded to the same precision.
func equalStats(a, b *Stats) error {
	if a.Len != b.Len {
		return fmt.Errorf("a.Len %v != b.Len %v", a.Len, b.Len)
	}
	if fmt.Sprintf("%.4f", a.Mean) != fmt.Sprintf("%.4f", b.Mean) {
		return fmt.Errorf("a.Mean %.4f != b.Mean %.4f", a.Mean, b.Mean)
	}
	if fmt.Sprintf("%.4f", a.StdDev) != fmt.Sprintf("%.4f", b.StdDev) {
		return fmt.Errorf("a.StdDev %.4f != b.StdDev %.4f", a.StdDev, b.StdDev)
	}
	if fmt.Sprintf("%.4f", a.Variance) != fmt.Sprintf("%.4f", b.Variance) {
		return fmt.Errorf("a.Variance %.4f != b.Variance %.4f", a.Variance, b.Variance)
	}
	if fmt.Sprintf("%.4f", a.Minimum) != fmt.Sprintf("%.4f", b.Minimum) {
		return fmt.Errorf("a.Minimum %.4f != b.Minimum %.4f", a.Minimum, b.Minimum)
	}
	if fmt.Sprintf("%.4f", a.Maximum) != fmt.Sprintf("%.4f", b.Maximum) {
		return fmt.Errorf("a.Maximum %.4f != b.Maximum %.4f", a.Maximum, b.Maximum)
	}
	if fmt.Sprintf("%.4f", a.Percentile25) != fmt.Sprintf("%.4f", b.Percentile25) {
		return fmt.Errorf("a.Percentile25 %.4f != b.Percentile25 %.4f", a.Percentile25, b.Percentile25)
	}
	if fmt.Sprintf("%.4f", a.Percentile50) != fmt.Sprintf("%.4f", b.Percentile50) {
		return fmt.Errorf("a.Percentile50 %.4f != b.Percentile50 %.4f", a.Percentile50, b.Percentile50)
	}
	if fmt.Sprintf("%.4f", a.Percentile75) != fmt.Sprintf("%.4f", b.Percentile75) {
		return fmt.Errorf("a.Percentile75 %.4f != b.Percentile75 %.4f", a.Percentile75, b.Percentile75)
	}
	if fmt.Sprintf("%.4f", a.Percentile95) != fmt.Sprintf("%.4f", b.Percentile95) {
		return fmt.Errorf("a.Percentile95 %.4f != b.Percentile95 %.4f", a.Percentile95, b.Percentile95)
	}
	if fmt.Sprintf("%.4f", a.Percentile99) != fmt.Sprintf("%.4f", b.Percentile99) {
		return fmt.Errorf("a.Percentile99 %.4f != b.Percentile99 %.4f", a.Percentile99, b.Percentile99)
	}
	return nil
}
