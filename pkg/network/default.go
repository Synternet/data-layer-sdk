package network

var defaultNetworks = Networks{
	defaultNetwork: "stub",
	all: map[string]Network{
		// Stub for testing. Does not establish actual network connection.
		"stub": {
			urls: []string{},
		},
		// Synternet testnet network.
		"testnet": {
			urls: []string{
				"europe-west3-gcp-dl-testnet-brokernode-frankfurt01.syntropynet.com",
				"europe-west3-gcp-dl-testnet-brokernode-frankfurt02.syntropynet.com",
				"europe-west3-gcp-dl-testnet-brokernode-frankfurt03.syntropynet.com",
			},
		},
	},
}
