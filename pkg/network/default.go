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
				"europe-west3-gcp-dl-testnet-brokernode-frankfurt01.synternet.com",
				"europe-west3-gcp-dl-testnet-brokernode-frankfurt02.synternet.com",
				"europe-west3-gcp-dl-testnet-brokernode-frankfurt03.synternet.com",
			},
		},
		// Synternet mainnet network.
		"mainnet": {
			urls: []string{
				"broker-eu-01.synternet.com",
				"broker-eu-02.synternet.com",
				"broker-eu-03.synternet.com",
			},
		},
	},
}
