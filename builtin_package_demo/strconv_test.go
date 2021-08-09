package builtin_package_demo

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

// int --> string
func TestItoa(t *testing.T) {
	i := 123
	s := strconv.Itoa(i)
	fmt.Println(s, reflect.TypeOf(s))
}
