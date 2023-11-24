package main

import (
	"bufio"
	"os"
	"silly/transport/proxy"
)

//
func Start() {
}

func main() {
	sock5 := proxy.Sock5Proxy{}
	sock5.Start(":8080", 0)
	stdin := bufio.NewReader(os.Stdin)
	stdin.ReadLine()
}
