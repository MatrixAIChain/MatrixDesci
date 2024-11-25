package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func GenerateKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		return nil, nil
	}
	return privateKey, &privateKey.PublicKey
}

func SignTransaction(privateKey *ecdsa.PrivateKey, data []byte) (r, s *big.Int, err error) {
	hash := sha256.Sum256(data)
	r, s, err = ecdsa.Sign(rand.Reader, privateKey, hash[:])
	return
}

func VerifySignature(publicKey *ecdsa.PublicKey, data []byte, r, s *big.Int) bool {
	hash := sha256.Sum256(data)
	return ecdsa.Verify(publicKey, hash[:], r, s)
}

func PublicKeyToAddress(pubKey *ecdsa.PublicKey) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%x", pubKey)))
	address := fmt.Sprintf("MRX-%s", hex.EncodeToString(hash[:]))[:34] // Address starts with MRX
	return address
}
