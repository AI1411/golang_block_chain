package message

import (
	"block_chain_go/pkg/protocol/common"
	"bytes"
	"fmt"
)

type Inv struct {
	Count     *common.VarInt
	Inventory []*common.InvVect
}

func NewInv(count *common.VarInt, inventory []*common.InvVect) *Inv {
	return &Inv{
		Count:     count,
		Inventory: inventory,
	}
}

func DecodeInv(b []byte) (*Inv, error) {
	inventory := []*common.InvVect{}
	varint, err := common.DecodeVarInt(b)
	if err != nil {
		return nil, err
	}
	length := len(varint.Encode())
	if uint64(len(b[length:])) != uint64(common.InventoryVectorSize)*varint.Data {
		return nil, fmt.Errorf("Decode to Inv failed, invalid input: %v", b)
	}
	b = b[length:]
	for i := 0; uint64(i) < varint.Data; i++ {
		invvect, err := common.DecodeInvVect(b[i*common.InventoryVectorSize : (i+1)*common.InventoryVectorSize])
		if err != nil {
			return nil, err
		}
		inventory = append(inventory, invvect)
	}
	return &Inv{
		Count:     varint,
		Inventory: inventory,
	}, nil
}


func (inv *Inv) CommandName() string {
	return "inv"
}

func (inv *Inv) Encode() []byte {
	inventoryBytes := [][]byte{}
	for _, invvect := range inv.Inventory {
		inventoryBytes = append(inventoryBytes, invvect.Encode())
	}
	return bytes.Join([][]byte{
		inv.Count.Encode(),
		bytes.Join(inventoryBytes, []byte{}),
	}, []byte{})
}