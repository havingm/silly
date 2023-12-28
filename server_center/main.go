package main

import (
	"bufio"
	"os"
	"silly/logger"
	"silly/silly"
	"silly/utils"
)

func main() {
	center, err := silly.NewCenter("127.0.0.1:8080", 0)
	if err != nil {
		logger.Error(err)
	}
	center.PrintHarbors()
	utils.SetConsoleTitle("Center[HarborId:0]")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadLine()
}
