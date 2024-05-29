package network

import (
	"fmt"
	"strings"
)

type Networks struct {
	Default string
	All     map[string]*Network
}

// SetDefault changes default network.
func SetDefault(name string) error {
	_, exists := defaultNetworks.All[name]

	if !exists {
		return fmt.Errorf("network '%s' does not exist", name)
	}

	defaultNetworks.Default = name

	return nil
}

// Default returns a default network.
func Default() *Network {
	network := defaultNetworks.All[defaultNetworks.Default]

	return network
}

type Network struct {
	URLs []string
}

func (n *Network) JoinURLs() string {
	urls := strings.Join(n.URLs, ",")

	return urls
}
