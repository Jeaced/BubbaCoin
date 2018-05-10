package main

import (
	"bubba_coin/core"
	"bubba_coin/cli"
)

func main() {
	bc := core.NewBlockchain()
	defer bc.CloseDb()

	cli := cli.NewCLI(bc)
	cli.Run()
}
