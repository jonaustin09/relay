package handler

import (
	"crypto/ecdsa"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/decred/base58"
	"github.com/getzion/relay/api/dwn/errors"
	. "github.com/getzion/relay/utils"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
)

// Get the public key from the attestation DID.
func (c *RequestContext) GetPublicKey() (*ecdsa.PublicKey, *errors.MessageLevelError) {
	// Validate the attestation object
	if c.Message.Attestation == nil {
		return nil, errors.NewMessageLevelError(400, "attestation cannot be null or empty", nil)
	} else if c.Message.Attestation.Protected == nil {
		return nil, errors.NewMessageLevelError(400, "attestation protected cannot be null", nil)
	} else if c.Message.Attestation.Protected.Alg != "ES256K" {
		return nil, errors.NewMessageLevelError(400, "Unsupported signing algorithm", nil)
	} else if strings.Trim(c.Message.Attestation.Protected.Kid, " ") == "" {
		return nil, errors.NewMessageLevelError(400, "Unsupported DID method", nil)
	} else if strings.HasPrefix(c.Message.Attestation.Protected.Kid, "did:key:") == false {
		return nil, errors.NewMessageLevelError(400, "Unsupported DID method", nil)
	}

	// Convert the multibase fingerprint to bytes.
	// This reverses the client-side method `getMultibaseFingerprintFromPublicKeyBytes`.
	fingerprintWithoutPrefix := strings.TrimPrefix(c.Message.Attestation.Protected.Kid, "did:key:z")
	didBytes := base58.Decode(fingerprintWithoutPrefix)
	pubKeyBytes := make([]byte, 33)
	pubKeyBytes = didBytes[2:]

	// Parse the pubkey bytes to verify that it corresponds to a valid public key on the secp256k1 curve.
	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return nil, errors.NewMessageLevelError(400, "Invalid pubkey", nil)
	} else {
		Log.Info().
			Bool("compressed", btcec.IsCompressedPubKey(pubKeyBytes)).
			Msg("Received valid pubkey")
	}

	// Convert the secp256k1 key to an ECDSA key.
	ecdsaKey := pubKey.ToECDSA()

	return ecdsaKey, nil
}

// Validate that this signature was generated by the given key.
func (c *RequestContext) VerifyRequest(publicKey *ecdsa.PublicKey) (bool, *errors.MessageLevelError) {
	signature := []byte(c.Message.Attestation.Signature)
	_, err3 := jws.Verify(signature, jwa.ES256K, publicKey)
	if err3 != nil {
		return false, nil
	} else {
		return true, nil
	}
}
