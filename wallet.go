package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"log"
	"bytes"
	"crypto/elliptic"
	"crypto/rand"
	"golang.org/x/crypto/ripemd160"
	"math/big"
	"fmt"
	"gx/ipfs/QmVmDhyTTUcQXFD1rRQ64fGLMSAoaQvNH3hwuaCFAPq2hy/errors"
	"encoding/hex"
)

const version = byte(0x00)
const walletFile = "wallet.dat"
const addressChecksumLen = 4

//Wallet stores private and public keys
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

//NewWallet creates and returns a Wallet
func NewWallet() *Wallet {
	private, public := newKeyPair()
	//private, public := SelfDefineKeyPair("0000000000000000000000000000000000000000000000000000000000000000")
	wallet := Wallet{private, public}

	return &wallet
}

//GetAddress returns wallet address

func (w Wallet) GetAddress() []byte  {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

//HashPubKey hashes public key
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return  publicRIPEMD160
}

//ValidateAddress check if address if valid
func ValidateAddress(address string) bool  {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash) - addressChecksumLen :]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

func checksum(payload []byte) []byte  {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

var (
	secp256k1_N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1_halfN = new(big.Int).Div(secp256k1_N, big.NewInt(2))
)

func toECDSA(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P256()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(secp256k1_N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return priv, nil
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


func SelfDefineKeyPair(hexkey string) (ecdsa.PrivateKey, []byte)  {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return ecdsa.PrivateKey{}, nil
	}
	pkey, err := toECDSA(b,true)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(pkey.PublicKey.X.Bytes(), pkey.PublicKey.Y.Bytes()...)

	//pubKeyHash := HashPubKey(pubKey)
	//
	//versionedPayload := append([]byte{version}, pubKeyHash...)
	//checksum := checksum(versionedPayload)
	//
	//fullPayload := append(versionedPayload, checksum...)
	//address := Base58Encode(fullPayload)

	return *pkey, pubKey
}