package builtin_package_demo

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

// TestGraceShutdown 优雅终止 http 服务
func TestGraceShutdown(t *testing.T) {
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

	done := make(chan bool)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	log.Println("Server is ready to handle requests at : 8080")
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on 8080: %v\n", err)
	}

	<-done
	log.Println("Server stopped")
}

func TestDumpRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%q", dump)
	}))
	defer ts.Close()

	const body = "test body"
	req, err := http.NewRequest("POST", ts.URL, strings.NewReader(body))
	if err != nil {
		t.Fatalf(err.Error())
	}
	req.Host = "www.example.com"
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println(string(b))
}
