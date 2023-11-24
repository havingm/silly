package main

import (
	"bufio"
	"os"
	"silly/transport/proxy"
)

func main() {
	server := new(proxy.TcpProxy)
	server.Start(":8080", ":6110")
	stdin := bufio.NewReader(os.Stdin)
	stdin.ReadLine()
}
