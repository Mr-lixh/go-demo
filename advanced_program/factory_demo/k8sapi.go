package factory_demo

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
)

type K8sApiType string

const (
	DefaultApiType K8sApiType = "default"
	Imp01ApiType   K8sApiType = "imp01"
)

// K8sApi is an interface implemented by things that know how
// to access kubernetes workloads.
type K8sApi interface {
	GetDeployment(namespace, name string) (*v1.Deployment, error)
	UpdateDeployment(deployment *v1.Deployment) error
	// other methods ...
}

func New(apiType K8sApiType) (K8sApi, error) {
	switch apiType {
	case Imp01ApiType:
		return NewImp01K8sApi()
	case DefaultApiType:
		return NewDefaultK8sApi()
	default:
		return nil, fmt.Errorf("invalid k8s api type: %v", apiType)
	}
}
