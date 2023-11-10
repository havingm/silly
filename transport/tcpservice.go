package transport

import (
	"fmt"
	"net"
	"silly/logger"
	"sync"
	"time"
)

type ITcpSocket interface {
	OnLinkRecv(link *TcpLink, data []byte)
	OnLinkClose(link *TcpLink)
}

type TcpService struct {
	sync.Mutex
	tag       string
	listener  *net.TCPListener
	linkTable sync.Map
	deadline  time.Duration
	close     chan struct{}
}

func NewTcpService(tag string) *TcpService {
	service := &TcpService{tag: tag}
	return service
}

func (s *TcpService) Start(addr string, deadline time.Duration) error {
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

func (s *TcpService) Stop() {
	if s.close == nil {
		return
	}
	select {
	case <-s.close:
		//已经关闭了
		return
	default:
		close(s.close)
	}
}

func (s *TcpService) OnLinkRecv(link *TcpLink, data []byte) {
	logger.Info("OnLinkRecv, linkId:", link.linkId, " 客户端发来: ", string(data))
	link.Send(data)
}

func (s *TcpService) OnLinkClose(link *TcpLink) {
	s.linkTable.Delete(link.linkId)
	logger.Info("OnLinkClose, linkId: ", link.linkId)
}

func (s *TcpService) CreateTcpLink(conn net.Conn) {
	s.Lock()
	s.Unlock()
	link := NewTcpLink(conn, s, 10)
	link.readDeadline = s.deadline
	link.linkId = AutoTcpLinkId()
	s.linkTable.Store(link.linkId, link)
	link.Run()

}

func (s *TcpService) accept() error {
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
	s.CreateTcpLink(conn)
	return nil
}

func (s *TcpService) run() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(r)
				logger.Error("Tcp服务发生错误退出")
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