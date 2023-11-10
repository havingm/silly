package gate_logic

import (
	"silly/transport"
)

type Gate struct {
	//用于客户端链接的Service
	clientService *transport.TcpService
	//用于服务端链接的Service
	serverService *transport.TcpService
	//
	clientManager *ClientManager
}
