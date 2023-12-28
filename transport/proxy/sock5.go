package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"silly/logger"
	"strconv"
	"sync"
	"time"
)

const socks5Version = 5

const (
	socks5AuthNone     = 0x00
	socks5AuthPassword = 0x02
)

const (
	socks5IP4    = 1
	socks5Domain = 3
	socks5Ip6    = 4
)

const (
	socks5CmdConnect = 0x01
	socks5CmdBind    = 0x02
	socks5CmdUdp     = 0x03
)

type Sock5Proxy struct {
	sync.Mutex
	tag      string
	listener *net.TCPListener
	deadline time.Duration
	close    chan struct{}
}

func (s *Sock5Proxy) Start(addr string, deadline time.Duration) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	s.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	s.deadline = deadline
	s.close = make(chan struct{})
	s.run()
	logger.Info(fmt.Sprintf("TcpService【%v】serv at addr: %v", s.tag, tcpAddr))
	return nil
}

/*
+----+----------+----------+
|VER | NMETHODS | METHODS |
+----+----------+----------+
| 1 | 1 | 1 to 255 |
+----+----------+----------+
*/
func (s *Sock5Proxy) handShake(conn net.Conn, reader *bufio.Reader) bool {
	if conn == nil || reader == nil {
		return false
	}
	b, err := reader.ReadByte()
	if err != nil {
		return false
	}
	//VER版本号，版本号必须为5
	if b != socks5Version {
		return false
	}
	//NMETHODS，METHODS的长度
	mLen, err := reader.ReadByte()
	if err != nil {
		return false
	}
	method := make([]byte, mLen)
	n, err := reader.Read(method)
	if n != int(mLen) || err != nil {
		return false
	}
	//返回无认证
	n, err = conn.Write([]byte{0x05, socks5AuthNone})
	if n != 2 || err != nil {
		return false
	}
	return true
}

/*
+----+-----+-------+------+----------+----------+
|VER | CMD | RSV | ATYP | DST.ADDR | DST.PORT |
+----+-----+-------+------+----------+----------+
| 1 | 1 | X'00' | 1 | Variable | 2 |
+----+-----+-------+------+----------+----------+
*/
func (s *Sock5Proxy) getAddress(reader *bufio.Reader) (host string, port string) {
	//先读固定的4个字节
	b := make([]byte, 4)
	n, err := reader.Read(b)
	if n != 4 || err != nil {
		return
	}
	cmd := ""
	switch b[1] {
	case socks5CmdConnect:
		cmd = "socks5CmdConnect"
	case socks5CmdBind:
		cmd = "socks5CmdBind"
	case socks5CmdUdp:
		cmd = "socks5CmdUdp"
	default:
		cmd = "socks5CmdUnknown"
	}
	logger.Info("[cmd:", cmd, "]")
	if b[0] != socks5Version ||
		(b[1] != socks5CmdConnect && b[1] != socks5CmdBind && b[1] != socks5CmdUdp) ||
		b[2] != 0x00 {
		return
	}
	var nport int
	typ := b[3]
	switch typ {
	case socks5IP4:
		/*
			the address is a version-4 IP address, with a length of 4 octets
		*/
		b = make([]byte, 6)
		n, err = reader.Read(b)
		if n != 6 || err != nil {
			return
		}
		ipv4 := b[:4]
		host = net.IP(ipv4).String()
		nport = int(b[4])<<8 + int(b[5])
	case socks5Domain:
		/*
			the address field contains a fully-qualified domain name. The first
			octet of the address field contains the number of octets of name that
			follow, there is no terminating NUL octet.
		*/
		len, err := reader.ReadByte()
		if len <= 0 || err != nil {
			return
		}
		b = make([]byte, len+2)
		n, err = reader.Read(b)
		if n != int(len+2) || err != nil {
			return
		}
		host = string(b[:len])
		nport = int(b[len])<<8 + int(b[len+1])
	case socks5Ip6:
		/*
			the address is a version-6 IP address, with a length of 16 octets.
		*/
		b = make([]byte, 18)
		n, err = reader.Read(b)
		if n != 18 || err != nil {
			return
		}
		ipv6 := b[:16]
		host = net.IP(ipv6).String()
		nport = int(b[16])<<8 + int(b[17])
	}
	return host, strconv.Itoa(nport)
}

func (s *Sock5Proxy) serveSock5(conn net.Conn) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(r)
			}
		}()
		reader := bufio.NewReader(conn)
		ok := s.handShake(conn, reader)
		if !ok {
			conn.Close()
			logger.Info("sock5 握手失败，关闭连接")
			return
		}
		logger.Info("sock5 握手成功")
		host, port := s.getAddress(reader)
		if len(host) <= 0 || port == "0" {
			conn.Close()
			logger.Info("sock5 获取地址失败，关闭连接")
			return
		}
		logger.Info("sock5 获取地址成功，host: ", host, " port: ", port)
		target, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			logger.Info("连接目标服务器失败：", err)
			conn.Close()
			return
		}
		defer conn.Close()
		defer target.Close()
		//todo
		conn.Write([]byte{5, 0, 0, 1, 0, 0, 0, 0, 0, 0})
		go func() {
			io.Copy(target, conn)
		}()
		io.Copy(conn, target)
	}()
}

func (s *Sock5Proxy) accept() error {
	conn, err := s.listener.AcceptTCP()
	if err != nil {
		return err
	}
	err = conn.SetKeepAlive(true)
	if err != nil {
		return err
	}
	err = conn.SetKeepAlivePeriod(time.Second * 30)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("TcpService【%v】accepted from : %v", s.tag, conn.RemoteAddr()))
	s.serveSock5(conn)
	return nil
}

func (s *Sock5Proxy) run() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(r)
			}
		}()
		for {
			err := s.accept()
			if err != nil {
				logger.Error(err)
				select {
				case <-s.close:
					return
				}
			}
		}
	}()
}
