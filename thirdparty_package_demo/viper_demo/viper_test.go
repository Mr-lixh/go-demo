package viper_demo

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"testing"
)

func TestLoadTomlFromIOReader(t *testing.T) {
	var tomlExample = []byte(`
[server]
host = "0.0.0.0"
port = 8088

[mysql]
host = "127.0.0.1"
port = "3306"
`)
	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(bytes.NewBuffer(tomlExample)); err != nil {
		t.Fatal(err)
	}
	fmt.Println(v.Get("server").(map[string]interface{})["host"])
}

func TestLoadTomlFromFile(t *testing.T) {
	v := viper.New()
	v.SetConfigFile("./testdata/test1.toml")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			t.Fatal("config file not found")
		} else {
			t.Fatal(err)
		}
	}

	fmt.Println(v.AllSettings())
}

func TestWatchTomlFromFile(t *testing.T) {
	v := viper.New()
	v.SetConfigFile("./testdata/test1.toml")
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed: ", e.Name)
	})

	select {}
}

func TestAccessNestedKeys(t *testing.T) {
	v := viper.New()
	v.SetConfigFile("./testdata/test1.toml")
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	if v.GetInt64("mysql.port") != int64(3306) {
		t.Fatal(fmt.Sprintf("get nested key value is %v, expected is 3306", v.GetInt64("mysql.port")))
	}
}

func TestExtractSubTree(t *testing.T) {
	v := viper.New()
	v.SetConfigFile("./testdata/test1.toml")
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	mysqlConfig := v.Sub("mysql")

	if mysqlConfig == nil {
		t.Fatal("mysql config not found")
	}

	fmt.Println(mysqlConfig.AllSettings())
}

func TestMigrateTomlFromFiles(t *testing.T) {
	v := viper.New()
	v.SetConfigFile("./testdata/test1.toml")
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	v.SetConfigFile("./testdata/test2.toml")
	if err := v.MergeInConfig(); err != nil {
		t.Fatal(err)
	}

	if v.GetStringMap("mysql")["port"].(int64) != 3317 {
		t.Fatal(fmt.Sprintf("mysql port is %v, expected is 3317\n", v.GetStringMap("mysql")["port"]))
	}
}

func TestMigrateTomlFromIOReader(t *testing.T) {
	v := viper.New()
	v.SetConfigFile("./testdata/test1.toml")
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	override := map[string]interface{}{
		"mysql": map[string]interface{}{
			"port": int64(3317),
		},
	}

	if err := v.MergeConfigMap(override); err != nil {
		t.Fatal(err)
	}

	if v.GetStringMap("mysql")["port"].(int64) != 3317 {
		t.Fatal(fmt.Sprintf("mysql port is %v, expected is 3317\n", v.GetStringMap("mysql")["port"]))
	}
}
