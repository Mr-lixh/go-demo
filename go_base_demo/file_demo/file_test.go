package file_demo

import (
	"os"
	"testing"
)

func WriteFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	s := "hello world!"
	_, err = file.WriteString(s)
	return err
}

func TestWriteFile(t *testing.T) {
	filename := "test.logger"
	err := WriteFile(filename)
	if err != nil {
		t.Fatal(err)
	}
}
