package main

import (
	"bufio"
	"os"
	"silly/server_center/logic"
)

func main() {
	center_logic.Start()
	reader := bufio.NewReader(os.Stdin)
	reader.ReadLine()
}
