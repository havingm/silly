package silly

import (
	"bytes"
	"errors"
	"github.com/vmihailenco/msgpack/v5"
	"net"
	"silly/logger"
	"silly/transport"
	. "silly/utils"
	"time"
)

type CenterHandler struct {
	Harbor *Harbor
}

func (c *CenterHandler) Init(harbor *Harbor) {
	c.Harbor = harbor
}

func (c *CenterHandler) OnLinkRecv(link *transport.TcpLink, data []byte) {
	logger.Info("CenterHandler.OnLinkRecved")
	buf := bytes.NewBuffer(data)
	msgId := ReadInt32(buf)
	logger.Info("msgId: ", msgId)
	methodName := ReadString(buf)
	paramBuf := ReadBytes(buf)
	switch methodName {
	case "CenterOnSyncAllHarbors":
		harborInfos := make(map[int]*HarborClient)
		err := msgpack.Unmarshal(paramBuf, &harborInfos)
		if err != nil {
			logger.Error("CenterOnSyncAllHarbors, err: ", err)
			return
		}
		c.Harbor.CenterOnSyncAllHarbors(harborInfos)
	}
}

func (c *CenterHandler) OnLinkClose(link *transport.TcpLink) {
	logger.Error("Harbor: ", c.Harbor.Id, " 与中心服的链接断开了，进行重连...")
	for i := 0; i < 10; i++ {
		if c.Harbor != nil {
			err := c.Harbor.ConnectCenter()
			if err == nil {
				break
			}
			logger.Error("重连失败：err", err, " 3秒后第：", i+1, " 次尝试...")
			time.Sleep(3 * time.Second)
		}
	}
}

type Harbor struct {
	HarborClient
	remoteReg    RemoteFunReg
	AllHarbors   map[int]*HarborClient
	NamedHarbors map[string]*HarborClient
	centerLink   *transport.TcpLink
	service      *transport.TcpService
	centerAddr   string
	centerHandle *CenterHandler
}

func (h *Harbor) OnLinkOpened(service *transport.TcpService, link *transport.TcpLink) {

}

func (h *Harbor) OnLinkClosed(service *transport.TcpService, link *transport.TcpLink) {

}

func (h *Harbor) OnLinkRecved(service *transport.TcpService, link *transport.TcpLink, data []byte) {
	logger.Info("Harbor.OnLinkRecved")
	methodName, params, err := UnmarshalServerCallMsg(data, &h.remoteReg, h)
	if err != nil {
		logger.Info(err)
		return
	}
	ReflectCallMethod(h, methodName, params...)
}

func (h *Harbor) ConnectCenter() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", h.centerAddr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	h.centerLink = transport.NewTcpLink(conn, h.centerHandle, 3)
	h.centerLink.Run()
	//上报harbor信息到中心服
	return h.RegHarbor2Center()

}

func NewHarbor(id int, name, addr string, deadline time.Duration, centerAddr string) (*Harbor, error) {
	harbor := &Harbor{}
	harbor.Id = id
	harbor.Name = name
	harbor.Tag = 1
	harbor.Addr = addr
	harbor.centerHandle = &CenterHandler{}
	harbor.centerHandle.Init(harbor)
	harbor.centerAddr = centerAddr
	harbor.AllHarbors = make(map[int]*HarborClient)
	harbor.NamedHarbors = make(map[string]*HarborClient)
	harbor.service = transport.NewTcpService(transport.WithTag(harbor.Name), transport.WithHolder(harbor))
	err := harbor.service.Start(addr, 0)
	harbor.remoteReg.Init()
	harbor.remoteReg.RegisterRemoteFunc(harbor)
	if err != nil {
		return nil, err
	}
	return harbor, nil
}

func (h *Harbor) Start() error {
	err := h.ConnectCenter()
	if err != nil {
		return err
	}
	return nil
}

func (h *Harbor) RegHarbor2Center() error {
	err := h.CallHarborMethod(0, "RemoteOnRegHarbor", h.HarborClient)
	return err
}

func (h *Harbor) CallHarborMethod(harborId int, methodName string, params ...interface{}) error {
	var targetLink *transport.TcpLink
	if harborId == 0 {
		if h.centerLink == nil {
			return errors.New("center link is nil")
		}
		targetLink = h.centerLink
	} else {
		target := h.AllHarbors[harborId]
		if target == nil || target.Link == nil {
			return errors.New("harbor is not exist or link unavailable")
		}
		targetLink = target.Link
	}

	buf, err := EncodeServerMsg(methodName, params...)
	if err != nil {
		return err
	}
	return targetLink.Send(buf)
}

func (h *Harbor) RemoteCallHarborOnHello(src int, msg string) {
	logger.Info("Harbor.RemoteCallHarborOnHello from harbor: ", src, " msg: ", msg)
}

func (h *Harbor) SayHelloToOtherHarbor(targetHarbor int, msg string) {
	h.CallHarborMethod(targetHarbor, "RemoteCallHarborOnHello", h.Id, msg)
}

func (h *Harbor) CenterOnSyncAllHarbors(harborClients map[int]*HarborClient) {
	logger.Info("Harbor.CenterOnSyncAllHarbors")
	for _, harborClient := range harborClients {
		if harborClient.Id == h.Id {
			continue
		}
		harbor := h.AllHarbors[harborClient.Id]
		if harbor == nil {
			harbor = harborClient
			h.AllHarbors[harborClient.Id] = harbor
			err := harbor.Connect()
			if err != nil {
				logger.Error(err)
			}
		} else {
			//addr is changed, reconnect
			if harbor.Addr != harborClient.Addr {
				harbor.Link.Close()
				harbor = harborClient
				err := harbor.Connect()
				if err != nil {
					logger.Error(err)
				}
			}
		}
	}
}
