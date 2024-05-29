package config

import (
	"fmt"
	"strings"
)

type Networks struct {
	Default string
	All     map[string]*Network
}

// SetDefaultNetwork changes default network.
func SetDefaultNetwork(name string) error {
	_, exists := defaultNetworks.All[name]

	if !exists {
		return fmt.Errorf("network '%s' does not exist", name)
	}

	defaultNetworks.Default = name

	return nil
}

// DefaultNetwork returns a default network.
func DefaultNetwork() *Network {
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
