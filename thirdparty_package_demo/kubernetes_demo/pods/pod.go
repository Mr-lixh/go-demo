package pods

import (
	"bytes"
	"context"
	"github.com/Mr-lixh/go-demo/thirdparty_package_demo/kubernetes_demo/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"strings"
)

//GetWorkdirByPod returns the workdir of a running pod
func GetWorkdirByPod(ctx context.Context, p *corev1.Pod, container string, config *rest.Config, c *kubernetes.Clientset) (string, error) {
	//cmd := []string{"sh", "-c", "echo $PWD"}
	cmd := []string{"pwd"}
	return execCommandInPod(ctx, p, container, cmd, config, c)
}

func CheckIfBashIsAvailable(ctx context.Context, p *corev1.Pod, container string, config *rest.Config, c *kubernetes.Clientset) bool {
	cmd := []string{"bash", "--version"}
	_, err := execCommandInPod(ctx, p, container, cmd, config, c)
	return err == nil
}

func execCommandInPod(ctx context.Context, p *corev1.Pod, container string, cmd []string, config *rest.Config, c *kubernetes.Clientset) (string, error) {
	in := strings.NewReader("\n")
	var out bytes.Buffer

	err := utils.Exec(
		ctx,
		c,
		config,
		p.Namespace,
		p.Name,
		container,
		false,
		in,
		&out,
		os.Stderr,
		cmd,
	)

	if err != nil {
		log.Printf("failed to execute command: %s - %s", err, out.String())
		return "", err
	}

	result := strings.TrimSuffix(out.String(), "\n")
	return result, nil
}
