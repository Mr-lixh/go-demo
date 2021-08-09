package yaml_demo

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

// 将两个yaml文件进行合并
func TestMigrateYaml(t *testing.T) {
	masterYaml := `
someProperty: "someValue"
anotherProperty: "anotherValue"`

	overrideYaml := `
someProperty: "overriddenValue"`

	var master map[string]interface{}
	if err := yaml.Unmarshal([]byte(masterYaml), &master); err != nil {
		t.Fatal(err)
	}

	var override map[string]interface{}
	if err := yaml.Unmarshal([]byte(overrideYaml), &override); err != nil {
		t.Fatal(err)
	}

	for k, v := range override {
		master[k] = v
	}

	newYaml, err := yaml.Marshal(master)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("overrided yaml is:\n%s\n", string(newYaml))
}
