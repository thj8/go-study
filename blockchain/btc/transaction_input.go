package btc

import "bytes"

// TXInput represents a transaction input
type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

// UseKey checks whether the address initiated the transaction
func (in *TXInput) UseKey(pubKeyHash []byte) bool {
	lockingHash := HashPubkey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
