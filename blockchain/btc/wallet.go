package btc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressChecksumLen = 4
const walletFile = "wallet.dat"

// Wallet stores private and public keys
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewWallet creates and returns a wallet
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

// GetAddress returns wallet address
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubkey(w.PublicKey)

	versionPlay := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionPlay)

	fullPlayload := append(versionPlay, checksum...)

	address := Base58Encode(fullPlayload)
	return address
}

// Checksum generates a checksum for a public key
func checksum(playload []byte) []byte {
	firstSHA := sha256.Sum256(playload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

// HashPubkey hashed public key
func HashPubkey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

// ValidateAddress check if address if valid
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	hashLen := len(pubKeyHash) - addressChecksumLen
	actualChecksum := pubKeyHash[hashLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1:hashLen]

	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(targetChecksum, actualChecksum) == 0
}
