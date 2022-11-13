package plugin_demo

import (
	"fmt"
	"github.com/golang/glog"
	"io"
	"os"
	"sync"
)

// Factory is a function that returns an Interface.
// The config parameter provides an io.Reader handler to the factory in
// order to load specific configurations. If no configuration is provided
// the parameter is nil.
type Factory func(config io.Reader) (Interface, error)

var (
	providersMutex sync.Mutex
	providers      = make(map[string]Factory)
)

// RegisterProvider registers a Factory by name. This is expected to happen during
// app startup.
func RegisterProvider(name string, provider Factory) {
	providersMutex.Lock()
	defer providersMutex.Unlock()

	if _, found := providers[name]; found {
		glog.Fatalf("Provider %q was registered twice", name)
	}

	glog.V(1).Infof("Registered providers %q", name)
	providers[name] = provider
}

// IsProvider returns true if name corresponds to an already registered providers.
func IsProvider(name string) bool {
	providersMutex.Lock()
	defer providersMutex.Unlock()

	_, found := providers[name]
	return found
}

// GetProvider creates an instance of the named providers, or nil if the name is unknown.
// The error return is only used if the named providers was known but failed to initialize.
// The config parameter specifies the io.Reader handler of the configuration file for the
// providers, or nil for no configuration.
func GetProvider(name string, config io.Reader) (Interface, error) {
	providersMutex.Lock()
	defer providersMutex.Unlock()

	f, found := providers[name]
	if !found {
		return nil, nil
	}
	return f(config)
}

// InitProvider creates an instance of the named providers.
func InitProvider(name string, configFilePath string) (Interface, error) {
	var provider Interface
	var err error

	if name == "" {
		glog.Info("No providers specified.")
		return nil, nil
	}

	if configFilePath != "" {
		var config *os.File
		config, err = os.Open(configFilePath)
		if err != nil {
			glog.Fatalf("Couldn't open providers configuration %s: %#v", configFilePath, err)
		}

		defer config.Close()
		provider, err = GetProvider(name, config)
	} else {
		provider, err = GetProvider(name, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("could not init providers %q: %v", name, err)
	}
	if provider == nil {
		return nil, fmt.Errorf("unknown providers %q", name)
	}

	return provider, nil
}
