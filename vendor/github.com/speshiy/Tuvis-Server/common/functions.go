package common

import "math"

//RoundValue round value for two symbols
func RoundValue(value float32) float32 {
	result := float64(value)
	result = math.Round(result * 100)
	result = result / 100
	return float32(result)
}
