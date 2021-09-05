package factory_demo

import "testing"

func TestNew(t *testing.T) {
	examples := []string{"default", "imp01", "imp02"}

	var k8sApi K8sApi
	var err error
	for _, e := range examples {
		k8sApi, err = New(K8sApiType(e))
		if err != nil {
			t.Errorf(err.Error())
		}
		if k8sApi == nil {
			t.Errorf("get %q k8s api instance is nil", e)
		}
	}
}
