package util

import "math"

func CalculateFee(satoshiPerByte uint64, utxoCount uint64) uint64 {
	var baseTransactionSize = 8 + 1 + (4+4+1+95+4)*utxoCount + 1 + 32*2

	totalTransactionSize := baseTransactionSize
	virtualTransactionSize := math.Ceil((float64(baseTransactionSize)*3 + float64(totalTransactionSize)) / 4)
	return uint64(virtualTransactionSize) * satoshiPerByte
}

func CalculateFeeForSegwit(satoshiPerByte uint64, utxoCount uint64) uint64 {
	var baseTransactionSize = 10 + 1 + (32+4+1+23+4)*utxoCount + 1 + 32*2

	witnessSize := 1 + 73 + 1 + 33
	totalTransactionSize := baseTransactionSize + 1 + uint64(witnessSize)*utxoCount
	virtualTransactionSize := math.Ceil((float64(baseTransactionSize)*3 + float64(totalTransactionSize)) / 4)
	return uint64(virtualTransactionSize) * satoshiPerByte
}