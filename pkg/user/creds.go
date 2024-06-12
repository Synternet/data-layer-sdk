package creds

import (
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

func CreateCreds(accSeed []byte, opts ...Opt) (seed string, jwt string, err error) {
	accKeys, err := nkeys.FromSeed(accSeed)
	if err != nil {
		return "", "", fmt.Errorf("failed to get account keys: %w", err)
	}

	userKeys, err := nkeys.CreateUser()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate user keys: %w", err)
	}

	userSeed, err := userKeys.Seed()
	if err != nil {
		return "", "", fmt.Errorf("failed to get user seed: %w", err)
	}

	userClaims, err := createUserClaims(accKeys, userKeys)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate user claims: %w", err)
	}

	var payload interface{} = userClaims
	if len(opts) > 0 {
		payloadMap, err := convertToMap(userClaims)
		if err != nil {
			return "", "", err
		}

		for _, opt := range opts {
			if err := opt(payloadMap); err != nil {
				return "", "", err
			}
		}

		payload = payloadMap
	}

	userJwt, err := encodeJWT(accKeys, payload)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate user claims: %w", err)
	}

	return string(userSeed), userJwt, nil
}

func createUserClaims(accKeys, userKeys nkeys.KeyPair) (*jwt.UserClaims, error) {
	accountPubKey, err := accKeys.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("error getting public key: %w", err)
	}

	userPubKey, err := userKeys.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("cannot get user's public key: %w", err)
	}

	claims := jwt.NewUserClaims(userPubKey)
	claims.Issuer = accountPubKey
	claims.Type = jwt.UserClaim
	claims.GenericFields.Version = 2

	c := &claims.ClaimsData
	c.Issuer = accountPubKey
	c.IssuedAt = time.Now().UTC().Unix()
	c.ID = "" // to create a repeatable hash
	c.ID, err = hashClaimsData(c)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func encodeJWT(accKeys nkeys.KeyPair, payload interface{}) (string, error) {
	header := &jwt.Header{Type: jwt.TokenTypeJwt, Algorithm: jwt.AlgorithmNkey}

	h, err := serializeAndEncode(header)
	if err != nil {
		return "", err
	}

	body, err := serializeAndEncode(payload)
	if err != nil {
		return "", err
	}

	toSign := fmt.Sprintf("%s.%s", h, body)

	if header.Algorithm != jwt.AlgorithmNkey {
		return "", errors.New(header.Algorithm + " not supported to write jwtV2")
	}

	sig, err := accKeys.Sign([]byte(toSign))
	if err != nil {
		return "", err
	}
	eSig := base64.RawURLEncoding.EncodeToString(sig)
	return fmt.Sprintf("%s.%s", toSign, eSig), nil
}

func hashClaimsData(c *jwt.ClaimsData) (string, error) {
	j, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	h := sha512.New512_256()
	h.Write(j)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(h.Sum(nil)), nil
}

func convertToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, err
}

func serializeAndEncode(v interface{}) (string, error) {
	j, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(j), nil
}
