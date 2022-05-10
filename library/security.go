package library

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/libp2p/go-libp2p-core/crypto"
)

func GenerateLibP2PKey() (string, error) {
	prk, _, err := crypto.GenerateSecp256k1Key(rand.Reader)
	if err != nil {
		return "", err
	}
	private_key_bytes, err := crypto.MarshalPrivateKey(prk)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(private_key_bytes), nil
}

func DecodeLibP2PKey(str string) (crypto.PrivKey, error) {
	private_key_bytes, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPrivateKey(private_key_bytes)
}
