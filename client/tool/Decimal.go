package tool

import (
	"fmt"
	"strconv"
)

// 截取小数点后两位
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
