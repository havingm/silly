package main

import (
	"bufio"
	"fmt"
	"os"
	"silly/logger"
	"silly/silly"
	"silly/utils"
	"strconv"
	"strings"
	"time"
)

var running bool
var reader *bufio.Reader
var server *silly.Harbor

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
	case "hello":
		if len(cmds) <= 1 {
			logger.Error("没有参数")
			return
		}
		params := strings.SplitN(cmds[1], " ", 2)
		if len(params) < 2 {
			logger.Error("参数不足")
			return
		}
		targetHarbor, _ := strconv.Atoi(params[0])
		message := params[1]
		server.SayHelloToOtherHarbor(targetHarbor, message)
	case "test":
	}
}

func main() {
	harborId := 1
	harborName := fmt.Sprintf("harbor_%d", harborId)
	addr := "127.0.0.1:8081"
	if len(os.Args) >= 4 {
		harborId, _ = strconv.Atoi(os.Args[1])
		harborName = os.Args[2]
		addr = os.Args[3]
	}
	harbor, err := silly.NewHarbor(harborId, harborName, addr, 0, "127.0.0.1:8080")
	if err != nil {
		logger.Info(err)
		return
	}
	err = harbor.Start()
	for err != nil {
		logger.Info("server launch failed, will retry after 3 seconds: ", err)
		time.Sleep(time.Second * 3)
		err = harbor.Start()
	}
	if err != nil {
		logger.Info(err)
		return
	}
	logger.Info(harbor)
	server = harbor
	utils.SetConsoleTitle(fmt.Sprintf("%s[HarborId:%d]", server.Name, server.Id))
	StartCommand()
}
