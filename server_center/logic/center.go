package center_logic

import (
	"bufio"
	"os"
	"silly/logger"
	"silly/transport"
)

var (
	//用于服务端链接的Service
	serverService *transport.TcpService
	//客户端链接管理器
	serverManager *ServerManager
)

type Gate struct {
}

func Init() {

}

//
func Start() {
	serverManager = &ServerManager{}
	serverService = transport.NewTcpService(transport.WithTag("gate_client_service"), transport.WithHolder(serverManager))
	err := serverService.Start("localhost:8080", 0)
	if err != nil {
		logger.Error(err)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	reader.ReadLine()
}
