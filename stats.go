package stats

import (
	"errors"
	"math"
	"sort"

	"golang.org/x/exp/constraints"
)

// Number is a constraint that permits any integer or floating-point type
type Number interface {
	constraints.Integer | constraints.Float
}

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

// ErrEmpty is returned when attempting to calculate stats on an empty slice.
var ErrEmpty = errors.New("empty array")

// CalcStats computes statistics on a slice of numbers.  It works on any
// integer or float type and does not modify the input slice.  It returns
// ErrEmpty if data is empty.
func CalcStats[T Number](data []T) (Stats, error) {
	n := len(data)
	if n == 0 {
		return Stats{}, ErrEmpty
	}

	// Create a float64 copy for calculations to avoid modifying the input and
	// to handle all number types uniformly.
	floatData := make([]float64, n)
	var sum, sumSq float64
	for i, v := range data {
		val := float64(v)
		floatData[i] = val
		sum += val
		sumSq += val * val
	}

	// Sort the copy for percentile, min, and max calculations.
	sort.Float64s(floatData)

	var variance, stdDev float64
	if n > 1 {
		// Sample variance
		variance = (sumSq - sum*sum/float64(n)) / float64(n-1)
		stdDev = math.Sqrt(variance)
	}

	stats := Stats{
		Len:          n,
		Mean:         sum / float64(n),
		Variance:     variance,
		StdDev:       stdDev,
		Minimum:      floatData[0],
		Maximum:      floatData[n-1],
		Percentile25: percentile(floatData, 25),
		Percentile50: percentile(floatData, 50),
		Percentile75: percentile(floatData, 75),
		Percentile95: percentile(floatData, 95),
		Percentile99: percentile(floatData, 99),
	}

	return stats, nil
}

// percentile calculates the k-th percentile of a sorted slice of data.  Uses
// the NIST-recommended "Nearest-Rank" method.
func percentile(sorted []float64, k float64) float64 {
	if k < 0 || k > 100 {
		panic("percentile must be between 0 and 100")
	}

	n := float64(len(sorted))
	if n == 1 {
		return sorted[0]
	}

	// Calculate index using N-1 scaling
	idx := int(math.Ceil(k / 100 * (n - 1)))

	return sorted[idx]
}
