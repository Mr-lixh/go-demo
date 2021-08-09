package builtin_package_demo

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
)

type MyHandler struct {
	foo string
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.String() {
	case "/test":
		fmt.Println("test url")
	case "/metrics":
		fmt.Println("metrics url")
	default:
		fmt.Println(h.foo)
	}
}

func TestMyHandler(t *testing.T) {
	h := &MyHandler{
		foo: "hello world!",
	}

	s := &http.Server{
		Addr:           ":8080",
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatalln(s.ListenAndServe())
}
