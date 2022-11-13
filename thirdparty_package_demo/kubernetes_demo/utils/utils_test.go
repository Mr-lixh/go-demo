package utils

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGetLocal(t *testing.T) {
	cs, _, _, err := GetLocal()
	if err != nil {
		t.Fatal(err)
	}

	pods, err := cs.CoreV1().Pods("default").List(context.Background(), v1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range pods.Items {
		fmt.Println(p.Name)
	}
}
