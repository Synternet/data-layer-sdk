package network

import (
	"fmt"
	"strings"
)

type Networks struct {
	defaultNetwork string
	all            map[string]Network
}

// SetDefault changes default network.
func SetDefault(name string) error {
	_, exists := defaultNetworks.all[name]

	if !exists {
		return fmt.Errorf("network '%s' does not exist", name)
	}

	defaultNetworks.defaultNetwork = name

	return nil
}

// Add adds a network.
func Add(name string, network Network) error {
	_, exists := defaultNetworks.all[name]

	if exists {
		return fmt.Errorf("network '%s' already exists", name)
	}

	// TODO: NATS checks network URLs when connecting, but feedback can be given earlier - here.
	// This needs to be thought through since it duplicates logic of checking URLs.
	defaultNetworks.all[name] = network

	return nil
}

// Default returns a default network.
func Default() Network {
	network := defaultNetworks.all[defaultNetworks.defaultNetwork]

	return network
}

type Network struct {
	urls []string
}

// JoinUrls is a utility function to join URLs into a string.
func (n Network) Urls() string {
	urls := strings.Join(n.urls, ",")

	return urls
}
