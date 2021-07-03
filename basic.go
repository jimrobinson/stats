package stats

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"sort"
)

// Stats holds statistical values computed from an array of numbers
type Stats struct {
	Len          int     // total number of entries
	Mean         float64 // average value across all entries
	StdDev       float64 // standard deviation between entries
	Variance     float64 // variance between entries
	Minimum      float64 // minimum value among all entries
	Maximum      float64 // maximum value among all entries
	Percentile25 float64 // 25% of all entries were less than or equal to this value
	Percentile50 float64 // 50% of all entries were less than or equal to this value
	Percentile75 float64 // 75% of all entries were less than or equal to this value
	Percentile95 float64 // 95% of all entries were less than or equal to this value
	Percentile99 float64 // 99% of all entries were less than or equal to this value
}

// SlotsFloat64 is a []float64 that provides Stats()
type SlotsFloat64 []float64

// sortFloat64 is a slice of pointers to SlotsFloat64
type sortFloat64 []*float64

// EMPTY is the error returned by Stats when the slots input array is
// empty, implying that no useful statistics can be computed from it.
var EMPTY = errors.New("empty array")

// Stats will compute some basic statistics on the values in slots.
// If the slots slice is empty an error will be returned.  The values
// for stats.StdDev and stats.Variance will only be computed if the
// slots slice has at least two values.
func (slots SlotsFloat64) Stats() (stats *Stats, err error) {
	stats = new(Stats)

	stats.Len = len(slots)
	if stats.Len == 0 {
		err = EMPTY
		return
	}

	var entries, sum, sumsq float64
	entries = float64(stats.Len)

	sorted := make(sortFloat64, stats.Len)
	for i, v := range slots {
		sorted[i] = &(slots[i])
		sum += v
		sumsq += v * v
	}
	if entries > 1 {
		stats.Variance = (1 / (entries - 1)) * (sumsq - (1/entries)*sum*sum)
		stats.StdDev = math.Sqrt(stats.Variance)
	}

	sort.Sort(sorted)
	stats.Mean = sum / entries
	stats.Minimum = *(sorted[0])
	stats.Maximum = *(sorted[stats.Len-1])
	stats.Percentile25 = *(sorted[(stats.Len*25)/100])
	stats.Percentile50 = *(sorted[(stats.Len*50)/100])
	stats.Percentile75 = *(sorted[(stats.Len*75)/100])
	stats.Percentile95 = *(sorted[(stats.Len*95)/100])
	stats.Percentile99 = *(sorted[(stats.Len*99)/100])

	return
}

func (slots sortFloat64) Len() int {
	return len(slots)
}

func (slots sortFloat64) Less(i, j int) bool {
	return *(slots[i]) < *(slots[j])
}

func (slots sortFloat64) Swap(i, j int) {
	slots[i], slots[j] = slots[j], slots[i]
}

func (s *Stats) String() string {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "Entries  : %d\n", s.Len)
	fmt.Fprintf(buf, "Mean     : %6.4f\n", s.Mean)
	fmt.Fprintf(buf, "StdDev   : %6.4f\n", s.StdDev)
	fmt.Fprintf(buf, "Variance : %6.4f\n", s.Variance)
	fmt.Fprintf(buf, "Minimum  : %6.4f\n", s.Minimum)
	fmt.Fprintf(buf, "Maximum  : %6.4f\n", s.Maximum)
	fmt.Fprintf(buf, "Percentiles\n")
	fmt.Fprintf(buf, "    25th : %6.4f\n", s.Percentile25)
	fmt.Fprintf(buf, "    50th : %6.4f\n", s.Percentile50)
	fmt.Fprintf(buf, "    75th : %6.4f\n", s.Percentile75)
	fmt.Fprintf(buf, "    95th : %6.4f\n", s.Percentile95)
	fmt.Fprintf(buf, "    99th : %6.4f", s.Percentile99)

	return buf.String()
}
