package kubernetes_demo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

// 根据 PID 获取 pod 信息

// ID identifies a single container running in a Kubernetes Pod
type ID struct {
	Namespace     string
	PodName       string
	PodUID        string
	PodLabels     map[string]string
	ContainerID   string
	ContainerName string
}

// LookupPod looks up a process ID from the host PID namespace, returning its Kubernetes identity.
func LookupPod(pid int) (*ID, error) {
	// First find the Docker container ID.
	cid, err := LookupDockerContainerID(pid)
	if err != nil {
		return nil, err
	}

	// Look up the container ID in the local kubelet API.
	resp, err := http.Get(fmt.Sprintf("http://localhost:%v/pods", 10255))
	if err != nil {
		return nil, errors.WithMessage(err, "could not lookup container ID in kubelet API")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "could not read response from kubelet API")
	}
	var podInfo *podList
	if err := json.Unmarshal(body, &podInfo); err != nil {
		return nil, errors.WithMessage(err, "could not unmarshal response from kubelet API")
	}

	for _, item := range podInfo.Items {
		for _, status := range item.Status.ContainerStatuses {
			if status.ContainerID == "docker://"+cid {
				return &ID{
					Namespace:     item.Metadata.Namespace,
					PodName:       item.Metadata.Name,
					PodUID:        item.Metadata.UID,
					PodLabels:     item.Metadata.Labels,
					ContainerID:   cid,
					ContainerName: status.Name,
				}, nil
			}
		}
	}
	return nil, nil
}

// LookupDockerContainerID looks up a process ID from the host PID namespace,
// returning its Docker container ID.
func LookupDockerContainerID(pid int) (string, error) {
	f, err := os.Open(fmt.Sprintf("/proc/%d/cgroup", pid))
	if err != nil {
		// this is normal, it just means the PID no longer exists
		return "", nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := dockerPattern.FindStringSubmatch(line)
		if parts != nil {
			return parts[1], nil
		}
		parts = kubePattern.FindStringSubmatch(line)
		if parts != nil {
			return parts[1], nil
		}
	}
	return "", nil
}

var (
	kubePattern   = regexp.MustCompile(`\d+:.+:/kubepods/[^/]+/pod[^/]+/([0-9a-f]{64})`)
	dockerPattern = regexp.MustCompile(`\d+:.+:/docker/pod[^/]+/([0-9a-f]{64})`)
)

type podList struct {
	// We only care about namespace, serviceAccountName and containerID
	Metadata struct {
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Namespace string            `json:"namespace"`
			Name      string            `json:"name"`
			UID       string            `json:"uid"`
			Labels    map[string]string `json:"labels"`
		} `json:"metadata"`
		Spec struct {
			ServiceAccountName string `json:"serviceAccountName"`
		} `json:"spec"`
		Status struct {
			ContainerStatuses []struct {
				ContainerID string `json:"containerID"`
				Name        string `json:"name"`
			} `json:"containerStatuses"`
		} `json:"status"`
	} `json:"items"`
}
