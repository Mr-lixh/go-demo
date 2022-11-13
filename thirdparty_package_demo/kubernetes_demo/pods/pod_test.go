package pods

import (
	"context"
	"fmt"
	"github.com/Mr-lixh/go-demo/thirdparty_package_demo/kubernetes_demo/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGetWorkdirByPod(t *testing.T) {
	cs, cc, _, err := utils.GetLocal()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	podList, err := cs.CoreV1().Pods("default").List(ctx, v1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(podList.Items) == 0 {
		t.Fatal("pod list is empty")
	}

	pod := podList.Items[1]
	fmt.Printf("pod name: %s, container name: %s\n", pod.Name, pod.Spec.Containers[0].Name)

	workdir, err := GetWorkdirByPod(ctx, &pod, pod.Spec.Containers[0].Name, cc, cs)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(workdir)
}
