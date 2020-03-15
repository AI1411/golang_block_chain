package main

import (
	"block_chain_go/internal/wallet"
	"block_chain_go/pkg/client"
	"block_chain_go/pkg/protocol/common"
	"block_chain_go/pkg/protocol/message"
	"block_chain_go/pkg/util"
	"log"
	"time"
)

func main() {
	k := util.NewKey()
	k.GenerateKey()
	log.Printf("key: %+v, %+v,%+v", k.PrivateKey.D.Bytes(), k.PublicKey.X.Bytes(), k.PublicKey.Y.Bytes())
	c := client.NewClient("testnet-seed.bitcoin.jonasschnelli.ch:18333")
	defer c.Conn.Close()
	log.Printf("remote addr: %s", c.Conn.RemoteAddr().String())

	addrFrom := &common.NetworkAddress{
		Services: uint64(1),
		IP: [16]byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x7F, 0x00, 0x00, 0x01,
		},
		Port: 8333,
	}

	v := &message.Version{
		Version:     uint32(70015),
		Services:    uint64(1),
		Timestamp:   uint64(time.Now().Unix()),
		AddrRecv:    addrFrom,
		AddrFrom:    addrFrom,
		Nonce:       uint64(0),
		UserAgent:   common.NewVarStr([]byte("")),
		StartHeight: uint32(0),
		Relay:       false,
	}

	_, err := c.SendMessage(v)
	if err != nil {
		log.Fatal(err)
	}

	wallet := wallet.NewWallet(c)
	wallet.Handshake()

	if err := wallet.Handshake(); err != nil {
		log.Fatal(err)
	}
	log.Printf("finish")
}
