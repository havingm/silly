package transport

import (
	"fmt"
	"net"
	"silly/logger"
	"sync"
	"time"
)

type Option func(*Options)
type Options struct {
	Tag    string
	Holder ITcpService
}

func NewOptions(opts ...Option) Options {
	opt := Options{
		Tag:    "Default",
		Holder: nil,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func WithTag(tag string) Option {
	return func(o *Options) {
		o.Tag = tag
	}
}

func WithHolder(holder ITcpService) Option {
	return func(o *Options) {
		o.Holder = holder
	}
}

type ITcpSocket interface {
	OnLinkRecv(link *TcpLink, data []byte)
	OnLinkClose(link *TcpLink)
}

type ITcpService interface {
	OnLinkOpened(service *TcpService, link *TcpLink)
	OnLinkClosed(service *TcpService, link *TcpLink)
	OnLinkRecv(service *TcpService, link *TcpLink, data []byte)
}

type TcpService struct {
	sync.Mutex
	tag      string
	holder   ITcpService
	listener *net.TCPListener
	deadline time.Duration
	close    chan struct{}
}

func NewTcpService(opts ...Option) *TcpService {
	options := NewOptions(opts...)
	service := &TcpService{
		tag:    options.Tag,
		holder: options.Holder,
	}
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
	if s.holder != nil {
		s.holder.OnLinkRecv(s, link, data)
	}
}

func (s *TcpService) OnLinkClose(link *TcpLink) {
	if s.holder != nil {
		s.holder.OnLinkClosed(s, link)
	}
}

func (s *TcpService) CreateTcpLink(conn net.Conn) {
	s.Lock()
	s.Unlock()
	link := NewTcpLink(conn, s, 10)
	link.readDeadline = s.deadline
	link.linkId = AutoTcpLinkId()
	link.Run()
	if s.holder != nil {
		s.holder.OnLinkOpened(s, link)
	}
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
func (s *TcpService) String() string {
	return fmt.Sprintf("[TcpService: %v]", s.tag)
}
