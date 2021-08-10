package viper_demo

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
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
	fmt.Println(v.GetStringMapString("mysql")["host"])
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

func TestMergeTomlFromFiles(t *testing.T) {
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

func TestMergeTomlFromIOReader(t *testing.T) {
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

func TestWriteToFile(t *testing.T) {
	const TempFile = "./testdata/test.toml"
	defer os.Remove(TempFile)

	example := []byte(`
[test]
host = "127.0.0.1"
port = 3306
`)
	v := viper.New()
	v.SetConfigFile(TempFile)
	if err := v.ReadConfig(bytes.NewBuffer(example)); err != nil {
		t.Fatal(err)
	}

	if err := v.WriteConfig(); err != nil {
		t.Fatal(err)
	}

	_, err := os.Stat(TempFile)
	if err != nil {
		t.Fatal(err)
	}

	otherV := viper.New()
	otherV.SetConfigFile(TempFile)
	if err = otherV.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	if otherV.GetStringMapString("test")["host"] != "127.0.0.1" {
		t.Fatal(fmt.Sprintf("get test.host is: %v, expected is: %v", otherV.GetStringMapString("test")["host"], "127.0.0.1"))
	}
}

func TestBindPFlag(t *testing.T) {
	pflag.Int("age", 11, "age of you")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		t.Fatal(err)
	}

	if viper.GetInt("age") != 11 {
		t.Fatal(fmt.Sprintf("get %v, expected %v", viper.GetInt("age"), 11))
	}
}

func TestUnmarshal(t *testing.T) {
	example := []byte(`
host = "127.0.0.1"
port = 3306
`)

	type Config struct {
		Host string
		Port int
	}

	var c Config

	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(bytes.NewBuffer(example)); err != nil {
		t.Fatal(err)
	}

	if err := v.Unmarshal(&c); err != nil {
		t.Fatal(err)
	}

	if c.Host != "127.0.0.1" {
		t.Fatal(fmt.Sprintf("get %v, expected %v", c.Host, "127.0.0.1"))
	}
}

func TestUnmarshalEmbedStruct(t *testing.T) {
	example := []byte(`
module:
    enabled: true
    token: 89h3f98hbwf987h3f98wenf89ehf
`)

	type moduleConfig struct {
		Token string
	}

	type config struct {
		Module struct {
			Enabled      bool
			moduleConfig `mapstructure:",squash"`
		}
	}

	var c config

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewBuffer(example)); err != nil {
		t.Fatal(err)
	}

	if err := v.Unmarshal(&c); err != nil {
		t.Fatal(err)
	}

	if c.Module.Token != "89h3f98hbwf987h3f98wenf89ehf" {
		t.Fatal(fmt.Sprintf("get %v", c.Module.Token))
	}
}

func TestMarshalToml(t *testing.T) {
	v := viper.New()
	v.SetConfigFile("./testdata/test1.toml")
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	c := v.AllSettings()
	bs, err := toml.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(bs))
}
