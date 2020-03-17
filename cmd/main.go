package main

import (
	"block_chain_go/internal/wallet"
	"block_chain_go/pkg/client"
	"block_chain_go/pkg/protocol/common"
	"block_chain_go/pkg/protocol/message"
	"bytes"
	"log"
)

func main() {
	c := client.NewClient("testnet-seed.bitcoin.jonasschnelli.ch:18333")
	defer c.Conn.Close()
	log.Printf("remote addr: %s", c.Conn.RemoteAddr().String())

	wallet := wallet.NewWallet(c)
	wallet.Handshake()

	if err := wallet.Handshake(); err != nil {
		log.Fatal(err)
	}

	pubkey := bytes.Join([][]byte{wallet.Key.PublicKey.X.Bytes(), wallet.Key.PublicKey.Y.Bytes()}, []byte{})
	wallet.Client.SendMessage(message.NewFilterload(1024, 10, [][]byte{pubkey}))

	size := uint32(common.MessageLen)
	for {
		buf ,_ := wallet.Client.ReceiveMessage(size)
		var header [24]byte
		copy(header[:], buf)
		msg := common.DecodeMessageHeader(header)
		log.Printf("receive msg.Langth: %+v", msg.Length)
		if bytes.HasPrefix(msg.Command[:], []byte("verack")) {
			log.Printf("receive verack")
			wallet.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("version")) {
			log.Printf("receive version")
			wallet.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("sendheaders")) {
			log.Printf("receive sendheaders")
			wallet.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("sendcmpct")) {
			log.Printf("receive sendcmpct")
			wallet.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("ping")) {
			log.Printf("receicve ping")
			b,_ := wallet.Client.ReceiveMessage(msg.Length)
			ping := message.DecodePing(b)
			pong := message.Pong{Nonce:ping.Nonce}
			wallet.Client.SendMessage(&pong)
		} else if bytes.HasPrefix(msg.Command[:], []byte("addr")) {
			log.Printf("receive addr")
			wallet.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("getheaders")) {
			log.Printf("receive getheaders")
			wallet.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("feefilter")) {
			log.Printf("recieve feefilter")
			wallet.Client.ReceiveMessage(msg.Length)
		} else {
			log.Printf("receive other")
			wallet.Client.ReceiveMessage(msg.Length)
		}
	}
	log.Printf("finish")
}
