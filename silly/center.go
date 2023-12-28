package silly

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"silly/logger"
	"silly/transport"
	. "silly/utils"
	"time"
)

type Center struct {
	Harbors map[int]*HarborClient
	service *transport.TcpService
}

func (c *Center) OnLinkOpened(service *transport.TcpService, link *transport.TcpLink) {
	logger.Info("Center.OnLinkOpened")
}

func (c *Center) OnLinkClosed(service *transport.TcpService, link *transport.TcpLink) {
	logger.Info("Center.OnLinkClosed")
	for _, harbor := range c.Harbors {
		if link == harbor.Link {
			c.Harbors[harbor.Id] = nil
		}
	}
}

func (c *Center) OnLinkRecved(service *transport.TcpService, link *transport.TcpLink, data []byte) {
	logger.Info("Center.OnLinkRecved")
	buf := bytes.NewBuffer(data)
	msgId := ReadInt32(buf)
	logger.Info("msgId: ", msgId)
	methodName := ReadString(buf)
	paramBuf := ReadBytes(buf)
	if methodName == "RemoteOnRegHarbor" {
		harborInfo := &HarborClient{}
		err := msgpack.Unmarshal(paramBuf, harborInfo)
		if err != nil {
			logger.Error("RemoteOnRegHarbor, err: ", err)
			return
		}
		harborInfo.Link = link
		c.RemoteOnRegHarbor(harborInfo)
	}
}

func (c *Center) CallHarborMethod(harborId int, methodName string, params ...interface{}) error {
	target := c.Harbors[harborId]
	if target == nil || target.Link == nil {
		return errors.New(fmt.Sprintf("harbor not exist or no harbor link, harborId: %d", harborId))
	}
	buf, err := EncodeServerMsg(methodName, params...)
	if err != nil {
		return err
	}
	return target.Link.Send(buf)
}

func (c *Center) RemoteOnRegHarbor(harborInfo *HarborClient) {
	c.Harbors[harborInfo.Id] = harborInfo
	c.PrintHarbors()
	c.SyncHarborTable2AllHarbor()
}

func (c *Center) SyncHarborTable2AllHarbor() {
	for _, harbor := range c.Harbors {
		err := c.CallHarborMethod(harbor.Id, "CenterOnSyncAllHarbors", c.Harbors)
		if err != nil {
			logger.Error("CallHarborMethod, CenterOnSyncAllHarbors err: ", err)
		}
	}
}

func (c *Center) PrintHarbors() {
	sInfo := "-------------------Harbor Tables Began------------------------\r\n"
	for _, hInfo := range c.Harbors {
		sInfo += hInfo.String()
	}
	sInfo += "-------------------Harbor Tables Ended------------------------\n"
	logger.Info(sInfo)
}

func NewCenter(addr string, deadline time.Duration) (*Center, error) {
	center := &Center{}
	center.Harbors = make(map[int]*HarborClient)
	center.service = transport.NewTcpService(transport.WithTag("Center"), transport.WithHolder(center))
	err := center.service.Start(addr, deadline)
	if err != nil {
		return nil, err
	}
	return center, nil
}
