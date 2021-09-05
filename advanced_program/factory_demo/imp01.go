package factory_demo

import v1 "k8s.io/api/apps/v1"

type Imp01 struct {
	// meta ...
}

func NewImp01K8sApi() (*Imp01, error) {
	return nil, nil
}

func (i *Imp01) GetDeployment(namespace, name string) (*v1.Deployment, error) {

	return nil, nil
}

func (i *Imp01) UpdateDeployment(d *v1.Deployment) error {

	return nil
}
