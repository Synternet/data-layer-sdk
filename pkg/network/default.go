package network

var defaultNetworks = Networks{
	Default: "stub",
	All: map[string]*Network{
		// Stub for testing. Does not establish actual network connection.
		"stub": {
			URLs: []string{},
		},
		// Synternet testnet network.
		"testnet": {
			URLs: []string{
				"europe-west3-gcp-dl-testnet-brokernode-frankfurt01.syntropynet.com",
				"europe-west3-gcp-dl-testnet-brokernode-frankfurt02.syntropynet.com",
				"europe-west3-gcp-dl-testnet-brokernode-frankfurt03.syntropynet.com",
			},
		},
	},
}
