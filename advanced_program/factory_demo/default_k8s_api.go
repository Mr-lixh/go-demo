package factory_demo

import (
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type k8s struct {
	client kubernetes.Interface
	// other meta ...
}

func NewDefaultK8sApi() (K8sApi, error) {
	client, err := kubernetes.NewForConfig(&rest.Config{})
	if err != nil {
		return nil, err
	}
	da := &k8s{
		client: client,
	}
	return da, nil
}

func (k *k8s) GetDeployment(namespace, name string) (*v1.Deployment, error) {
	return k.client.AppsV1().Deployments(namespace).Get(name, v12.GetOptions{})
}

func (k *k8s) UpdateDeployment(d *v1.Deployment) error {

	return nil
}
