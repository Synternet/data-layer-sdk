package config

import "fmt"

type Networks struct {
	Default string
	All     map[string]*Network
}

type Network struct {
	URLs []string
}

// SetDefault changes default network.
func (n *Networks) SetDefault(name string) error {
	_, exists := n.All[name]

	if !exists {
		return fmt.Errorf("network '%s' does not exist", name)
	}

	n.Default = name

	return nil
}

// GetDefault returns a default network.
func (n *Networks) GetDefault() Network {
	network := n.All[n.Default]

	return *network
}
