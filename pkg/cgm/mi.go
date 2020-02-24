package cgm

import "math"

// Absi is a int variant of math.Abs which is float64
func Absi(a int) int {
	return int(math.Abs(float64(a)))
}

// Mini is a int variant of math.Min which is float64
func Mini(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

// Maxi is a int variant of math.Max which is float64
func Maxi(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

// Sqrti is a int variant of math.Sqrt which is float64
func Sqrti(a int) int {
	return int(math.Sqrt(float64(a)))
}
