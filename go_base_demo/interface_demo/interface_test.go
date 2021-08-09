package interface_demo

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type Writer interface {
	Write([]byte) (int, error)
}

func WriteString(s string, w Writer) error {
	b := []byte(s)
	_, err := w.Write(b)
	return err
}

func TestWriterInterface(t *testing.T) {
	f, _ := os.Open(os.DevNull)
	var s strings.Builder

	_ = WriteString("hello world!", f)
	_ = WriteString("hello world!", &s)

	fmt.Println(s.String())
}
