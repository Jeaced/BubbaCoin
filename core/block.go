package core

import (
	"time"
	"bytes"
	"encoding/gob"
)

type Block struct {
	Header *BlockHeader
	Data   *BlockData
}

type BlockHeader struct {
	HashComplexity    int
	Timestamp         int64
	PreviousBlockHash []byte
	Hash              []byte
	Nonce             int
}

type BlockData struct {
	Payload []byte
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}

	return result.Bytes()

}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}

	return &block
}

func NewBlock(data string, previousBlockHash []byte, hashComplexity int) *Block {
	block := &Block{&BlockHeader{hashComplexity, time.Now().Unix(),
		previousBlockHash, nil, 0}, &BlockData{[]byte(data)}}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Header.Hash = hash
	block.Header.Nonce = nonce

	return block
}