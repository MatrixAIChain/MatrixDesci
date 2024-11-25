package transaction

import (
	"crypto/ecdsa"
	"fmt"
	"matrix-blockchain/utils"
	"time"
)

type Transaction struct {
	From      string
	To        string
	Amount    int64
	Timestamp int64
	Signature *ecdsa.Signature
}

// NewTransaction creates a new transaction, signs it with the sender's private key
func NewTransaction(from, to string, amount int64, privateKey *ecdsa.PrivateKey) (*Transaction, error) {
	// Create a new transaction with necessary details
	transaction := &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}

	// Signing the transaction by hashing the details and using the private key
	hash := fmt.Sprintf("%s:%d:%d", transaction.From, transaction.Amount, transaction.Timestamp)
	r, s, err := utils.SignTransaction(privateKey, []byte(hash))
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}
	transaction.Signature = &ecdsa.Signature{R: r, S: s}

	return transaction, nil
}

// Verify checks the transaction's signature to ensure its authenticity
func (t *Transaction) Verify() bool {
	// Hash the transaction details (from, to, amount, timestamp)
	hash := fmt.Sprintf("%s:%d:%d", t.From, t.Amount, t.Timestamp)

	// Retrieve the sender's public key from their address and verify the signature
	publicKey := utils.GetPublicKeyFromAddress(t.From) // This should map an address to its public key
	return utils.VerifySignature(publicKey, []byte(hash), t.Signature.R, t.Signature.S)
}
