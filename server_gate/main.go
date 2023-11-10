package main

import (
	"bufio"
	"os"
	gate_logic "silly/server_gate/logic"
)

func main() {
	gate_logic.Start()
	reader := bufio.NewReader(os.Stdin)
	reader.ReadLine()
}
