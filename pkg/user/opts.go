package user

import "fmt"

type Opt func(map[string]interface{}) error

const (
	keyNats             = "nats"
	keyJetStreamManager = "jetstream_manager"
)

func JetStreamManagerOpt(m map[string]interface{}) error {
	natsUser, ok := m[keyNats]
	if !ok {
		return fmt.Errorf("failed to inject JetStream manager option, as there is no nats setup")
	}

	natsUserMap, ok := natsUser.(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to inject JetStream manager option, as nats is not a map value")
	}
	natsUserMap[keyJetStreamManager] = true
	return nil
}
