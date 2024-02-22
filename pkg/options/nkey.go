package options

import (
	"fmt"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

// CreateUser creates NATS user NKey and JWT from given account seed NKey.
func CreateUser(seed string) (*string, *string, error) {
	accountSeed := []byte(seed)

	accountKeys, err := nkeys.FromSeed(accountSeed)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get account key from seed: %w", err)
	}

	accountPubKey, err := accountKeys.PublicKey()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting public key: %w", err)
	}

	userKeys, err := nkeys.CreateUser()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create account key: %w", err)
	}

	userSeed, err := userKeys.Seed()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get seed: %w", err)
	}
	nkey := string(userSeed)

	userPubKey, err := userKeys.PublicKey()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot get user's public key: %w", err)
	}

	claims := jwt.NewUserClaims(userPubKey)
	claims.Issuer = accountPubKey
	jwt, err := claims.Encode(accountKeys)
	if err != nil {
		return nil, nil, fmt.Errorf("error encoding token to jwt: %w", err)
	}

	return &nkey, &jwt, nil
}
