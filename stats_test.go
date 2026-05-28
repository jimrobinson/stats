package stats

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

// epsilon is the small constant used to decide whether or not two floats are
// equal to one another
const epsilon = 1e-4

func TestStatsF64(t *testing.T) {
	statTests := []struct {
		Slots    []float64
		Expected Stats
	}{
		{
			Slots: []float64{1, 1},
			Expected: Stats{
				Len:          2,
				Mean:         1,
				StdDev:       0,
				Variance:     0,
				Minimum:      1,
				Maximum:      1,
				Percentile25: 1,
				Percentile50: 1,
				Percentile75: 1,
				Percentile95: 1,
				Percentile99: 1}},
		{
			Slots: []float64{2, 1},
			Expected: Stats{
				Len:          2,
				Mean:         1.5,
				StdDev:       0.7071,
				Variance:     0.5,
				Minimum:      1,
				Maximum:      2,
				Percentile25: 2,
				Percentile50: 2,
				Percentile75: 2,
				Percentile95: 2,
				Percentile99: 2}},
		{
			Slots: []float64{4, 3, 2, 1},
			Expected: Stats{Len: 4,
				Mean:         2.5,
				StdDev:       1.291,
				Variance:     1.6667,
				Minimum:      1,
				Maximum:      4,
				Percentile25: 2,
				Percentile50: 3,
				Percentile75: 4,
				Percentile95: 4,
				Percentile99: 4}},
		{
			Slots: []float64{256, 1, 1024, 512},
			Expected: Stats{Len: 4,
				Mean:         448.25,
				StdDev:       436.8618,
				Variance:     190848.25,
				Minimum:      1,
				Maximum:      1024,
				Percentile25: 256,
				Percentile50: 512,
				Percentile75: 1024,
				Percentile95: 1024,
				Percentile99: 1024}},
		{
			Slots: []float64{2, 1, 5, 4, 3, 6, 9, 8, 7, 10},
			Expected: Stats{Len: 10,
				Mean:         5.5,
				StdDev:       3.0277,
				Variance:     9.1667,
				Minimum:      1,
				Maximum:      10,
				Percentile25: 4,
				Percentile50: 6,
				Percentile75: 8,
				Percentile95: 10,
				Percentile99: 10}},
	}

	// test zero length array
	stats, err := CalcStats([]float64{})
	if err != ErrEmpty {
		t.Errorf("expected error value of %v, got %v", ErrEmpty, err)
	}

	// test valid arrays
	for i := range statTests {
		stats, err = CalcStats(statTests[i].Slots)
		if err != nil {
			t.Errorf("unexpected error value on statTests[%d].Slots.Stats(): %v", i, err)
		}
		err = equalStats(&statTests[i].Expected, &stats)
		if err != nil {
			t.Logf("Slots=%v", statTests[i].Slots)
			t.Logf("Expect=%#v", statTests[i].Expected)
			t.Logf("Actual=%#v", stats)
			t.Error(err)
		}
	}
}

func BenchmarkStats1K(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		slots := make([]float64, 1<<10)
		for i := range slots {
			slots[i] = rand.Float64()
		}
		b.StartTimer()
		CalcStats(slots)
		b.StopTimer()
	}
}

func BenchmarkStats64K(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		slots := make([]float64, 1<<16)
		for i := range slots {
			slots[i] = rand.Float64()
		}
		b.StartTimer()
		CalcStats(slots)
		b.StopTimer()
	}
}

// floatsApproxEqual verifies whether the difference between two numbers is
// smaller than a defined constant epsilon
func floatsApproxEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

// equalStats examines Stats a and b, and returns an error if any of the
// members do not match when rounded to the same precision.
func equalStats(a, b *Stats) error {
	if a.Len != b.Len {
		return fmt.Errorf("a.Len %v != b.Len %v", a.Len, b.Len)
	}
	if !floatsApproxEqual(a.Mean, b.Mean) {
		return fmt.Errorf("a.Mean %.4f != b.Mean %.4f", a.Mean, b.Mean)
	}
	if !floatsApproxEqual(a.StdDev, b.StdDev) {
		return fmt.Errorf("a.StdDev %.4f != b.StdDev %.4f", a.StdDev, b.StdDev)
	}
	if !floatsApproxEqual(a.Variance, b.Variance) {
		return fmt.Errorf("a.Variance %.4f != b.Variance %.4f", a.Variance, b.Variance)
	}
	if !floatsApproxEqual(a.Minimum, b.Minimum) {
		return fmt.Errorf("a.Minimum %.4f != b.Minimum %.4f", a.Minimum, b.Minimum)
	}
	if !floatsApproxEqual(a.Maximum, b.Maximum) {
		return fmt.Errorf("a.Maximum %.4f != b.Maximum %.4f", a.Maximum, b.Maximum)
	}
	if !floatsApproxEqual(a.Percentile25, b.Percentile25) {
		return fmt.Errorf("a.Percentile25 %.4f != b.Percentile25 %.4f", a.Percentile25, b.Percentile25)
	}
	if !floatsApproxEqual(a.Percentile50, b.Percentile50) {
		return fmt.Errorf("a.Percentile50 %.4f != b.Percentile50 %.4f", a.Percentile50, b.Percentile50)
	}
	if !floatsApproxEqual(a.Percentile75, b.Percentile75) {
		return fmt.Errorf("a.Percentile75 %.4f != b.Percentile75 %.4f", a.Percentile75, b.Percentile75)
	}
	if !floatsApproxEqual(a.Percentile95, b.Percentile95) {
		return fmt.Errorf("a.Percentile95 %.4f != b.Percentile95 %.4f", a.Percentile95, b.Percentile95)
	}
	if !floatsApproxEqual(a.Percentile99, b.Percentile99) {
		return fmt.Errorf("a.Percentile99 %.4f != b.Percentile99 %.4f", a.Percentile99, b.Percentile99)
	}
	return nil
}
