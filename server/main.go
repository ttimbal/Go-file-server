package main

import (
	. "file-server/cli"
	. "file-server/utils"
)

func main() {
	PrintReset()
	server := NewServer()
	cli := NewCli("server")
	cli.Register(server)
	go cli.Start()

	for active := range cli.Active {
		if !active {
			break
		}
	}

}
