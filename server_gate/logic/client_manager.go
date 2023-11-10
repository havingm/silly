package gate_logic

import "silly/transport"

type ClientManager struct {
}

func (mgr *ClientManager) OnLinkRecv(link *transport.TcpLink, data []byte) {
	panic("implement me")
}

func (mgr *ClientManager) OnLinkClose(link *transport.TcpLink) {
	panic("implement me")
}
