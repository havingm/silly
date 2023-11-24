package gate_logic

import (
	"silly/logger"
	"silly/transport"
	"sync"
)

type ClientManager struct {
	linkTable sync.Map
}

func (mgr *ClientManager) OnLinkOpened(service *transport.TcpService, link *transport.TcpLink) {
	logger.Info("OnLinkOpened: ", service, link)
	mgr.linkTable.Store(link.GetLinkId(), link)
}

func (mgr *ClientManager) OnLinkClosed(service *transport.TcpService, link *transport.TcpLink) {
	logger.Info("OnLinkClosed: ", service, link)
	mgr.linkTable.Delete(link.GetLinkId())
}

func (mgr *ClientManager) OnLinkRecved(service *transport.TcpService, link *transport.TcpLink, data []byte) {
	logger.Info("OnLinkRecv: ", service, link, data)
}
