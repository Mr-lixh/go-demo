package _default

import (
	"github.com/Mr-lixh/go-demo/advanced_program/plugin_demo"
	"io"
)

const (
	ProviderName = "default"
)

// Provider is an implementation of Interface.
type Provider struct {
	// some meta
	client string
	// ......
}

func (p *Provider) Method1() (string, error) {
	// imply Method1
	return "", nil
}

func (p *Provider) Method2() error {
	// imply Method2
	return nil
}

// newDefaultProvider creates a new instance of DefaultProvider.
func newDefaultProvider(config io.Reader) (*Provider, error) {
	// init providers config by config or use default config
	// ......

	return CreateDefaultProvider(config)
}

// CreateDefaultProvider creates a DefaultProvider object using the specified parameters.
func CreateDefaultProvider(config io.Reader) (*Provider, error) {
	p := &Provider{
		client: "fake-clent",
	}

	return p, nil
}

func init() {
	plugin_demo.RegisterProvider(
		ProviderName,
		func(config io.Reader) (plugin_demo.Interface, error) {
			return newDefaultProvider(config)
		})
}
