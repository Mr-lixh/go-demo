package _default

import (
	"github.com/Mr-lixh/go-demo/advanced_program/plugin_demo"
	"testing"
)

func TestInitProvider(t *testing.T) {
	var examples = []struct {
		Name string
	}{
		{
			Name: "default",
		},
		{
			Name: "test",
		},
	}
	for _, e := range examples {
		_, err := plugin_demo.InitProvider(e.Name, "")
		if err != nil {
			t.Errorf("Provider could not be initialized: %v", err)
			break
		}
		t.Logf("Provider %q has been registered", e.Name)
	}
}
