package message

import (
	"block_chain_go/pkg/protocol/common"
	"block_chain_go/pkg/util"
	"bytes"
	"encoding/binary"
	"github.com/spaolacci/murmur3"
	"math"
)

type Filterload struct {
	Count      *common.VarInt
	Filter     []byte
	NHashFuncs uint32
	NTweak     uint32
	NFlags     uint8
}

func NewFilterload(size uint32, nHashfuncs uint32, queries [][]byte) *Filterload {
	byteArray := make([]byte, size)
	nTweak := make([]byte, 4)
	for i := 0; i < cap(nTweak); i++ {
		nTweak[i] = util.RandInt8(0, math.MaxUint8)
	}
	nTweakUint32 := binary.BigEndian.Uint32(nTweak)
	for _, query := range queries {
		for i := 0; uint32(i) < nHashfuncs; i++ {
			seed := uint32(i)*0xFBA4C795 + nTweakUint32
			hashValue := murmur3.Sum32WithSeed(query, seed)
			adjustHashValue := hashValue % (size * uint32(8))
			idx := 1 << (uint32(7) & hashValue)
			value := 1 << (uint32(7) & hashValue)
			byteArray[idx] = byte(value)
		}
	}
	return &Filterload{
		Count:      common.NewVarInt(uint64(size)),
		Filter:     byteArray,
		NHashFuncs: nHashfuncs,
		NTweak:     nTweakUint32,
		NFlags:     uint8(1),
	}
}

func (f *Filterload) Command() [12]byte {
	var commandName [12]byte
	copy(commandName[:], "filterload")
	return commandName
}

func (f *Filterload) Encode() []byte {
	nHashFuncsBytes := make([]byte, 4)
	nTweakByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(nHashFuncsBytes, f.NHashFuncs)
	binary.LittleEndian.PutUint32(nTweakByte, f.NHashFuncs)
	return bytes.Join([][]byte{
		f.Count.Encode(),
		f.Filter,
		nHashFuncsBytes,
		nTweakByte,
		[]byte{f.NFlags},
	}, []byte{})
}