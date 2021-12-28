package tool

import "math"

//转换内存 uint64转换为float64

func MemTrans(value uint64, f int) float64 {
	return Decimal(float64(value)/math.Pow10(f), 0)
}
