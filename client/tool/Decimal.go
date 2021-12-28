package tool

import "math"

// 对于Float64类型 截取小数点后两位
func Decimal(value float64, fx int) float64 {
	f := math.Pow10(fx)
	return math.Trunc(value*f) / f
}
