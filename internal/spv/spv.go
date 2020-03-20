package spv

import (
	"block_chain_go/pkg/client"
	"block_chain_go/pkg/protocol/common"
	"block_chain_go/pkg/protocol/message"
	"block_chain_go/pkg/util"
	"bytes"
	"encoding/hex"
	"log"
)

type SPV struct {
	Client  *client.Client
	Key     *util.Key
	Address string
	Balance uint64
}

func NewSPV(client *client.Client) *SPV {
	key := util.NewKey()
	key.GenerateKey()
	serializedPubKey := key.PublicKey.SerializeUncompressed()
	address := util.EncodeAddress(serializedPubKey)
	//log.Printf("address:%s", address)
	return &SPV{
		Client:  client,
		Key:     key,
		Address: address,
		Balance: 0,
	}
}

func (s *SPV) Handshake() error {
	v := message.NewVersion()
	_, err := s.Client.SendMessage(v)
	if err != nil {
		return err
	}
	var recvVerack, recvVersion bool
	for {
		if recvVerack && recvVersion {
			log.Printf("success handshake")
			return nil
		}
		buf, err := s.Client.ReceiveMessage(common.MessageLen)
		if err != nil {
			log.Printf("handshake Receive massage error: %+v", err)
			return err
		}

		var header [24]byte
		copy(header[:], buf)
		msg := common.DecodeMessageHeader(header)
		_, err = s.Client.ReceiveMessage(msg.Length)
		if err != nil {
			return err
		}

		if bytes.HasPrefix(msg.Command[:], []byte("verack")) {
			recvVerack = true
		} else if bytes.HasPrefix(msg.Command[:], []byte("version")) {
			recvVersion = true
			_, err := s.Client.SendMessage(&message.Verack{})
			if err != nil {
				return err
			}
		}
	}
}

func (s *SPV) MessageHandler() {
	for {
		buf, err := s.Client.ReceiveMessage(common.MessageLen)
		if err != nil {
			//log.Printf("message handler err: %+v", err)
			log.Fatal("message handler err: ", err)
			//continue
		}
		var header [24]byte
		copy(header[:], buf)
		msg := common.DecodeMessageHeader(header)

		if bytes.HasPrefix(msg.Command[:], []byte("verack")) {
			s.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("version")) {
			s.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("sendheaders")) {
			s.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("sendcmpct")) {
			s.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("ping")) {
			b, _ := s.Client.ReceiveMessage(msg.Length)
			ping := message.DecodePing(b)
			pong := message.Pong{
				Nonce: ping.Nonce,
			}
			s.Client.SendMessage(&pong)
		} else if bytes.HasPrefix(msg.Command[:], []byte("addr")) {
			s.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("getheaders")) {
			s.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("feefilter")) {
			s.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("inv")) {
			log.Printf("msg: %+v", msg)
			b, _ := s.Client.ReceiveMessage(msg.Length)
			inv, _ := message.DecodeInv(b)
			log.Printf("inv: %+v", inv.Count)

			inventory := []*message.InvVect{}
			for _, iv := range inv.Inventory {
				if iv.Type == message.InvTypeMsgBlock {
					inventory = append(inventory, message.NewInvVect(message.InvTypeMsgFilteredBlock, iv.Hash))
				}
			}
			s.Client.SendMessage(message.NewGetData(inventory))
		} else if bytes.HasPrefix(msg.Command[:], []byte("merkleblock")) {
			b, _ := s.Client.ReceiveMessage(msg.Length)
			mb, _ := message.DecodeMerkleBlock(b)
			log.Printf("merkleblock: %+v", mb)
			h := mb.GetBlockHash()
			hexHash := hex.EncodeToString(util.ReverseBytes(h[:]))
			log.Printf("BlockHash: %s", hexHash)
			log.Printf("Hashes length:%+v", len(mb.Hashes))
			txHashes := mb.Validate()
			log.Printf("txhashes: %+v", txHashes)
			inventory := []*message.InvVect{}
			for _, txHash := range txHashes {
				inventory = append(inventory, message.NewInvVect(message.InvTypeMsgTx, txHash))
			}
			s.Client.SendMessage(message.NewGetData(inventory))
		} else if bytes.HasPrefix(msg.Command[:], []byte("tx")) {
			b, _ := s.Client.ReceiveMessage(msg.Length)
			tx, _ := message.DecodeTx(b)
			log.Printf("tx: %+v", tx)
			log.Printf("txID: %+v", tx.ID())
		} else {
			log.Printf("receive: other")
			s.Client.ReceiveMessage(msg.Length)
		}
	}
}
