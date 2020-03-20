package main

import (
	"block_chain_go/internal/wallet"
	"block_chain_go/pkg/client"
	"block_chain_go/pkg/protocol/message"
	"block_chain_go/pkg/util"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func main() {
	c := client.NewClient("testnet-seed.bitcoin.jonasschnelli.ch:18333")
	defer c.Conn.Close()
	log.Printf("remote addr: %s", c.Conn.RemoteAddr().String())

	wallet := wallet.NewWallet(c)
	wallet.Handshake()

	if err := wallet.Handshake(); err != nil {
		log.Fatal("handshake error", err)
	}

	pubkey := bytes.Join([][]byte{wallet.Key.PublicKey.X.Bytes(), wallet.Key.PublicKey.Y.Bytes()}, []byte{})
	wallet.Client.SendMessage(message.NewFilterload(1024, 10, [][]byte{pubkey}))

	startBlockHash, err := hex.DecodeString("0000000000000657bda6681e1a3d1aac92d09d31721e8eedbca98cac73e93226")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var arr [32]byte
	copy(arr[:], util.ReverseBytes(startBlockHash))
	getblocks := message.NewGetBlocks(uint32(700015), [][32]byte{arr}, message.ZeroHash)
	wallet.Client.SendMessage(getblocks)

	wallet.MessageHandler()
	log.Printf("finish")
}

