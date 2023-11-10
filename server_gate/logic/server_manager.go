package gate_logic

import "silly/transport"

type ServerManager struct {
}

func (mgr *ServerManager) OnLinkRecv(link *transport.TcpLink, data []byte) {
	panic("implement me")
}

func (mgr *ServerManager) OnLinkClose(link *transport.TcpLink) {
	panic("implement me")
}
