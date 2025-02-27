package bls

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/prysmaticlabs/prysm/v5/crypto/bls/blst"
	blscommon "github.com/prysmaticlabs/prysm/v5/crypto/bls/common"
	"os"
)

// We are using similar structure for saving bls keys as ethereum keystore
// https://github.com/ethereum/go-ethereum/blob/master/accounts/keystore/key.go
//
// We are storing PubKey sepearately so that we can list the pubkey without
// needing password to decrypt the private key
type EncryptedBLSKeyJSONV3 struct {
	PubKey string              `json:"pubKey"`
	Crypto keystore.CryptoJSON `json:"crypto"`
}

func ReadPrivateKeyFromFile(path, password string) (blscommon.SecretKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var encryptedBLSStruct EncryptedBLSKeyJSONV3
	err = json.Unmarshal(data, &encryptedBLSStruct)
	if err != nil {
		return nil, err
	}

	// Check if pubkey is present, if not return error
	// There is an issue where if you specify ecdsa key file
	// it still works and returns a keypair since the format of storage is same.
	// This is to prevent and make sure pubkey is present.
	// ecdsa keys doesn't have that field
	if encryptedBLSStruct.PubKey == "" {
		return nil, fmt.Errorf("invalid bls key file. pubkey field not found")
	}

	skBytes, err := keystore.DecryptDataV3(encryptedBLSStruct.Crypto, password)
	if err != nil {
		return nil, err
	}
	var key blscommon.SecretKey
	key, err = blst.SecretKeyFromBytes(skBytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}
