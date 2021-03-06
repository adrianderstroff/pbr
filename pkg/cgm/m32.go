package cgm

import "math"

// Abs32 is a float32 variant of math.Abs which is float64
func Abs32(a float32) float32 {
	return float32(math.Abs(float64(a)))
}

// Min32 is a float32 variant of math.Min which is float64
func Min32(a, b float32) float32 {
	return float32(math.Min(float64(a), float64(b)))
}

// Max32 is a float32 variant of math.Max which is float64
func Max32(a, b float32) float32 {
	return float32(math.Max(float64(a), float64(b)))
}

// Sqrt32 is a float32 variant of math.Sqrt which is float64
func Sqrt32(a float32) float32 {
	return float32(math.Sqrt(float64(a)))
}

// Floor32 is a float32 variant of math.Floor which is float64
func Floor32(a float32) float32 {
	return float32(math.Floor(float64(a)))
}

// Mod32 is a float32 variant of math.Mod which is float64
func Mod32(a, b float32) float32 {
	return float32(math.Mod(float64(a), float64(b)))
}

// Sin32 is a float32 variant of math.Sin which is float64
func Sin32(a float32) float32 {
	return float32(math.Sin(float64(a)))
}

// Cos32 is a float32 variant of math.Cos which is float64
func Cos32(a float32) float32 {
	return float32(math.Cos(float64(a)))
}

// Asin32 is a float32 variant of math.Asin which is float64
func Asin32(a float32) float32 {
	return float32(math.Asin(float64(a)))
}

// Atan232 is a float32 variant of math.Atan2 which is float64
func Atan232(a, b float32) float32 {
	return float32(math.Atan2(float64(a), float64(b)))
}
