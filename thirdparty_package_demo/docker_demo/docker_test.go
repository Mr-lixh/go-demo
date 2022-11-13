package docker_demo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"
	"testing"
)

// 拉取镜像并启动一个容器
func TestRunContainer(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 拉取镜像
	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	// 创建容器
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 启动容器
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		t.Fatalf(err.Error())
	}

	// 获取容器状态
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf(err.Error())
		}
	case <-statusCh:
	}

	// 获取容器标准输出日志
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		t.Fatalf(err.Error())
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}

// 查询本地所有的容器
func TestListContainers(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatalf(err.Error())
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, c := range containers {
		fmt.Println(c.ID)
	}
}

// 查询本地所有的镜像
func TestListImages(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatalf(err.Error())
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, i := range images {
		fmt.Println(i.ID)
	}
}

// 拉取镜像
func TestPullImage(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatalf(err.Error())
	}

	authConfig := types.AuthConfig{
		Username: "admin",
		Password: "Harbor12345",
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := cli.ImagePull(ctx, "alpine", types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer out.Close()

	io.Copy(os.Stdout, out)
}

// 获取容器进程信息
func TestTopContainer(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatalf(err.Error())
	}

	body, err := cli.ContainerTop(ctx, "26ccd503b2faf03a3311e16df88bc3956ef5714157500f7a902d608f0d8e2415", []string{})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 获取容器所有 PID
	var PIDs []string
	var index int
	for i, t := range body.Titles {
		if t == "PID" {
			index = i
			break
		}
	}
	for _, proc := range body.Processes {
		PIDs = append(PIDs, proc[index])
	}

	fmt.Println(PIDs)
}
