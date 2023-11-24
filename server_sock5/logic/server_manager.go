package server_sock5

import (
	"silly/logger"
	"silly/transport"
	"sync"
)

type ServerManager struct {
	linkTable sync.Map
}

func (mgr *ServerManager) OnLinkOpened(service *transport.TcpService, link *transport.TcpLink) {
	logger.Info("OnlinkOpened: ", service, link)
	mgr.linkTable.Store(link.GetLinkId(), link)
}

func (mgr *ServerManager) OnLinkClosed(service *transport.TcpService, link *transport.TcpLink) {
	logger.Info("OnLinkClosed: ", service, link)
	mgr.linkTable.Delete(link.GetLinkId())
}

func (mgr *ServerManager) OnLinkRecved(service *transport.TcpService, link *transport.TcpLink, data []byte) {
	logger.Info("OnLinkRecv: ", service, link, data)
}
