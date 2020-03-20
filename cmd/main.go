package main

import (
	"block_chain_go/internal/wallet"
	"block_chain_go/pkg/client"
	"block_chain_go/pkg/protocol/message"
	"block_chain_go/pkg/util"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func main() {
	c := client.NewClient("[2604:a880:2:d0::2065:5001]:18333")
	defer c.Conn.Close()
	log.Printf("remote addr: %s", c.Conn.RemoteAddr().String())
	wallet := wallet.NewWallet(c)
	wallet.Handshake()

	if err := wallet.Handshake(); err != nil {
		log.Fatal("handshake error", err)
	}

	pubKeyHash := util.Hash160(wallet.Key.PublicKey.SerializeUncompressed())
	wallet.Client.SendMessage(message.NewFilterload(1024, 10, [][]byte{pubKeyHash}))

	startBlockHash, err := hex.DecodeString("000000000000020c54ca0a429835b14ba2f1629562547d39a0523af5dd518865")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var reversedStartBlockHash [32]byte
	copy(reversedStartBlockHash[:], util.ReverseBytes(startBlockHash))
	getblocks := message.NewGetBlocks(uint32(70015), [][32]byte{reversedStartBlockHash}, message.ZeroHash)
	wallet.Client.SendMessage(getblocks)

	wallet.MessageHandler()
	log.Printf("finish")
}
