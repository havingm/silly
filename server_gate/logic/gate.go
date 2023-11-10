package gate_logic

import (
	"bufio"
	"os"
	"silly/logger"
	"silly/transport"
)

var (
	//用于客户端链接的Service
	clientService *transport.TcpService
	//用于服务端链接的Service
	serverService *transport.TcpService
	//客户端链接管理器
	clientManager *ClientManager
)

type Gate struct {
}

func Init() {

}

//
func Start() {
	clientManager = &ClientManager{}
	clientService = transport.NewTcpService(transport.WithTag("gate_client_service"), transport.WithHolder(clientManager))
	err := clientService.Start("localhost:8080", 0)
	if err != nil {
		logger.Error(err)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	reader.ReadLine()
}
