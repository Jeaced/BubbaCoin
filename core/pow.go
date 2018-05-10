package core

import (
	"math/big"
	"bytes"
	"bubba_coin/utils"
	"fmt"
	"math"
	"crypto/sha256"
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.Header.PreviousBlockHash,
			pow.block.Data.Payload,
			utils.IntToHex(pow.block.Header.Timestamp),
			utils.IntToHex(int64(pow.block.Header.HashComplexity)),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var (
		hashInt big.Int
		hash    [32]byte
	)

	nonce := 0
	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data.Payload)
	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Header.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-b.Header.HashComplexity))
	pow := &ProofOfWork{b, target}

	return pow
}
