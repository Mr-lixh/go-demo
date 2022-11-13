package num_demo

import (
	"fmt"
	"math"
	"testing"
)

// 计算两个整数的百分比，向上取整
func TestDevision(t *testing.T) {
	var a = 15
	var b = 16
	c := float64(a) / float64(b) * 100
	percent := int(math.Ceil(c))

	fmt.Printf("当前百分比：%d%%\n", percent)
}
