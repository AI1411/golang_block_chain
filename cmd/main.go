package main

import (
	"block_chain_go/internal/spv"
	"block_chain_go/pkg/client"
	"block_chain_go/pkg/protocol/message"
	"block_chain_go/pkg/util"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func main() {
	c := client.NewClient("[2001:41d0:a:f7eb::1]:18333")
	defer c.Conn.Close()
	log.Printf("remote addr: %s", c.Conn.RemoteAddr().String())
	spv := spv.NewSPV(c)
	spv.Handshake()

	if err := spv.Handshake(); err != nil {
		log.Fatal("handshake error", err)
	}

	pubKeyHash := util.Hash160(spv.Key.PublicKey.SerializeUncompressed())
	spv.Client.SendMessage(message.NewFilterload(1024, 10, [][]byte{pubKeyHash}))

	startBlockHash, err := hex.DecodeString("000000000000020c54ca0a429835b14ba2f1629562547d39a0523af5dd518865")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var reversedStartBlockHash [32]byte
	copy(reversedStartBlockHash[:], util.ReverseBytes(startBlockHash))
	getblocks := message.NewGetBlocks(uint32(70015), [][32]byte{reversedStartBlockHash}, message.ZeroHash)
	spv.Client.SendMessage(getblocks)

	if err := spv.MessageHandler(); err != nil {
		log.Printf("main: message handler err:", err)
	}

	log.Printf("finish")
}
