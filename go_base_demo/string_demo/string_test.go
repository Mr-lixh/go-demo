package string_demo

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestMultilineString(t *testing.T) {
	str := `
This is a
multiline
string.
  Contain white space.`
	fmt.Println(str)

	str = `
This string
	will have tabs in it.`
	fmt.Println(str)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var source = rand.NewSource(time.Now().UnixNano())

// GenerateRandomString 生成特定长度的随机字符串
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[source.Int63()%int64(len(charset))]
	}
	return string(b)
}

func TestGenerateRandomString(t *testing.T) {
	var s string

	for i := 3; i < 6; i++ {
		s = GenerateRandomString(i)
		fmt.Println(s)
	}
}
