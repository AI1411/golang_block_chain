package util

import (
	"bytes"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/bech32"
	"log"
)

func EncodeAddress(publicKeyBytes []byte) string {
	bs := bytes.Join([][]byte{
		[]byte{0x6F},
		Hash160(publicKeyBytes),
	}, []byte{})

	checksum := Hash256(bs)[:4]
	return base58.Encode(bytes.Join([][]byte{bs, checksum}, []byte{}))
}

func EncodeNativeSegwitAddress(publicKeyBytes []byte) string {
	bs := bytes.Join([][]byte{
		[]byte{0x6F},
		Hash160(publicKeyBytes),
	}, []byte{})

	b5, err := bech32.ConvertBits(bs, 8, 5, true)
	log.Printf("err: %+v", err)
	b,_ := bech32.Encode("tb", b5)
	return b
}
