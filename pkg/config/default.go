package config

var DefaultNetworks = Networks{
	Default: "stub",
	All: map[string]*Network{
		"stub": {
			URLs: []string{},
		},
		"testnet": {
			URLs: []string{
				"europe-west3-gcp-dl-testnet-brokernode-frankfurt01.syntropynet.com",
				"europe-west3-gcp-dl-testnet-brokernode-frankfurt02.syntropynet.com",
				"europe-west3-gcp-dl-testnet-brokernode-frankfurt03.syntropynet.com",
			},
		},
	},
}
