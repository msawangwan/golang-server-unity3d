package main

import (
	"github.com/msawangwan/unity-server/network"
)

func main() {
	server := network.NewServerInstance(":9081")
	server.Start()
	server.Shutdown()
}
