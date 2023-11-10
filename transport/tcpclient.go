package transport

import (
	"net"
)

type TcpClient struct {
	conn *net.TCPConn
	addr string
}

func (c *TcpClient) ConnectServer(addr string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	c.conn = tcpConn
	return nil
}

func (c *TcpClient) Send(data []byte) error {
	return nil
}

func (c *TcpClient) run() {

}
