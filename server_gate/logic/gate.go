package gate_logic

import (
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

//网关客户端对象
type GatePlayer struct {
	PlayerId uint64 //玩家Id(未登录时为0)
	ConnId   uint64 //链接Id
	AgentId  uint64 //对应的AgentId
}

type GateConn struct {
	ConnId   uint64
	PlayerId uint64
}

type Gate struct {
	Players     map[uint64]*GatePlayer
	Connections map[uint64]*GateConn
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
	StartCommand()
}
