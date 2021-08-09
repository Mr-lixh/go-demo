package builtin_package_demo

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"testing"
	"time"
)

func TestPprof(t *testing.T) {
	pprofIp := "0.0.0.0:8086"
	go func() {
		err := http.ListenAndServe(pprofIp, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}()

	var buf []byte
	tick := time.Tick(time.Second / 100)
	for range tick {
		buf = append(buf, make([]byte, 1024*1024)...)
	}
}
