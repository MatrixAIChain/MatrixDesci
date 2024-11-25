package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

// Serialize converts an interface to a byte slice using gob encoding.
func Serialize(data interface{}) []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

// Deserialize converts a byte slice back to an interface using gob decoding.
func Deserialize(data []byte, object interface{}) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(object)
	if err != nil {
		log.Panic(err)
	}
}

// Hash calculates the SHA-256 hash of a byte slice.
func Hash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// ValidateAddress checks if a given address starts with "MRX".
func ValidateAddress(address string) bool {
	return len(address) > 4 && address[:3] == "MRX"
}

// CheckError is a utility function to handle errors conveniently.
func CheckError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
