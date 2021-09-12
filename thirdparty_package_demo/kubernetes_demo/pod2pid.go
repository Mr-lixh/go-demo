package kubernetes_demo

type ContainerInfo struct {
	ContainerID string
	PIDs        []string
	Ports       []int
}

type ContainerInfoList []*ContainerInfo

func GetContainerInfosByPod(namespace, name string) (ContainerInfoList, error) {
	// 1、调用 kubelet api获取pod对应的容器id列表
	// 2、根据容器id获取容器详情
	return nil, nil
}

// GetContainerInfoById 根据容器 ID 获取容器的 PIDs 和 Ports
func GetContainerInfoById(id string) (*ContainerInfo, error) {
	// 1、调用docker top api获取容器内的PIDs
	// 2、调用netstat api获取进程对应的端口
	return nil, nil
}
