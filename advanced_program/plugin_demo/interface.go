package plugin_demo

// Interface is an abstract, pluggable interface for providers.
type Interface interface {
	Method1() (string, error)
	Method2() error
	// other methods ...
}
