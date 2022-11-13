package slice_demo

import (
	"fmt"
	"testing"
)

func TestCreateSlice(t *testing.T) {
	var x, y, z = 3, 2, 3

	var threeDimsSlice = make([][][]int, x)
	for i := range threeDimsSlice {
		threeDimsSlice[i] = make([][]int, y)
		for j := range threeDimsSlice[i] {
			threeDimsSlice[i][j] = make([]int, z)
		}
	}

	// 每个元素为初始值
	fmt.Println(threeDimsSlice)
}

func TestCopySlice(t *testing.T) {
	slice1 := []int{0, 1, 2, 3, 4, 5, 6, 7}
	fmt.Println(slice1)

	// 切片操作是引用操作，修改切片后的数组，也会导致原切片发生变化
	// 浅拷贝：不同切片共用同一底层数组
	var slice2 = make([]int, 0)
	slice2 = slice1[:]
	slice2[0] = 100
	fmt.Println(slice1)

	// 深拷贝：源切片和目标切片彼此有独立的底层数组空间，修改不影响到其它切片
	slice1 = []int{0, 1, 2, 3, 4, 5, 6, 7}
	slice3 := make([]int, len(slice1))
	copy(slice3, slice1)
	slice3[0] = 9
	fmt.Println(slice1)
	fmt.Println(slice3)
}
