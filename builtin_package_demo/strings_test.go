package builtin_package_demo

import (
	"fmt"
	"strings"
	"testing"
)

func TestStringBuilder(t *testing.T) {
	var b strings.Builder
	for i := 3; i >= 1; i-- {
		fmt.Fprintf(&b, "%d...", i)
	}

	b.Write([]byte("ready..."))

	b.WriteString("ignition!!!")
	fmt.Println(b.String())
	fmt.Printf("builder len: %d, cap: %d\n", b.Len(), b.Cap())

	// reset
	b.Reset()
	fmt.Printf("after reset, builder: %s, builder len: %d, cap: %d\n", b.String(), b.Len(), b.Cap())
}
