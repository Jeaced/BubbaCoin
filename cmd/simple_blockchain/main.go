package main

import (
	"simle_blockchain/core"
	"fmt"
	"strconv"
)

func main() {
	bc := core.NewBlockchain(16)
	defer bc.CloseDb()

	bc.AddBlock("Send 1 BubbaCoin to Jeaced")
	bc.AddBlock("Send 5 BubbaCoin to Jeaced")

	bci := bc.Iterator()

	for  {
		block := bci.Next()
		fmt.Printf("Prev. hash: %x\n", block.Header.PreviousBlockHash)
		fmt.Printf("Data: %s\n", block.Data.Payload)
		fmt.Printf("Hash: %x\n", block.Header.Hash)
		pow := core.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.Header.PreviousBlockHash) == 0 {
			break
		}
	}

}
