package gate_logic

import (
	"bufio"
	"os"
	"silly/logger"
	"strings"
)

var running bool
var reader *bufio.Reader

func StartCommand() {
	//go func() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
		}
	}()
	running = true
	reader = bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		DoCommand(string(data))
	}
	//}()
}

func DoCommand(command string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
		}
	}()
	if len(command) <= 0 {
		return
	}
	logger.Info("收到指令: ", command)
	cmds := strings.SplitN(command, " ", 2)
	if len(cmds) <= 0 {
		return
	}
	cmd := cmds[0]
	//params := ""
	//if len(cmds) > 1 {
	//	params = cmds[1]
	//}
	switch cmd {
	case "s":
		running = false
	}
}
