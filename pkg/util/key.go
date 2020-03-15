package util

import (
	"crypto/rand"
	"github.com/btcsuite/btcd/btcec"
	"io/ioutil"
	"os"
)

const keyPath = "./privatekey"

type Key struct {
	PrivateKey *btcec.PrivateKey
	PublicKey  *btcec.PublicKey
}

func NewKey() *Key {
	return &Key{}
}

func (k *Key) GenerateKey() error {
	if existsFile(keyPath) {
		randomBytes, err := ioutil.ReadFile(keyPath)
		if err != nil {
			return err
		}
		privateKey, publicKey := btcec.PrivKeyFromBytes(btcec.S256(), randomBytes)
		k.PrivateKey = privateKey
		k.PublicKey = publicKey
	} else {
		randomBytes, err := generateRandom()
		if err != nil {
			return err
		}
		writeFile(randomBytes)
		privateKey, publicKey := btcec.PrivKeyFromBytes(btcec.S256(), randomBytes)
		k.PrivateKey = privateKey
		k.PublicKey = publicKey
	}
	return nil
}

func generateRandom() ([]byte, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

func writeFile(b []byte) error {
	err := ioutil.WriteFile(keyPath, b,0666)
	if err != nil {
		return err
	}
	return nil
}

func existsFile(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return !os.IsNotExist(err)
	}

	return err == nil
}
