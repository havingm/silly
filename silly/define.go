package silly

import (
	"fmt"
	"net"
	"silly/transport"
)

type HarborClient struct {
	Id   int
	Addr string
	Tag  int
	Name string
	Link *transport.TcpLink
}

func (h *HarborClient) OnLinkRecv(link *transport.TcpLink, data []byte) {

}

func (h *HarborClient) OnLinkClose(link *transport.TcpLink) {

}

func (h *HarborClient) Connect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", h.Addr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	h.Link = transport.NewTcpLink(conn, h, 3)
	h.Link.Run()
	return nil
}

func (h *HarborClient) String() string {
	info := fmt.Sprintf("Harbor Id: %v, Addr: %v, Tag: %v, Name: %v \r\n", h.Id, h.Addr, h.Tag, h.Name)
	return info
}
