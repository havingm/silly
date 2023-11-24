package proxy

import (
	"errors"
	"io"
	"net"
	"silly/logger"
)

//转发代理服
type TcpProxy struct {
	addrTarget string       //转发地址
	listener   net.Listener //tcp监听对象
	closeChan  chan struct{}
}

//启动
func (s *TcpProxy) Start(serveAddr string, targetAddr string) error {
	if s.closeChan != nil {
		select {
		case <-s.closeChan:
			//已经关闭了
		default:
			return errors.New("重复启动")
		}
	}

	logger.Info("转发代理服启动 监听地址：", serveAddr, " 转发地址：", targetAddr)

	var err error
	s.listener, err = net.Listen("tcp", serveAddr)
	if err != nil {
		logger.Info("转发代理服监听失败", err)
		return err
	}

	s.addrTarget = targetAddr
	s.closeChan = make(chan struct{})

	s.Run()

	return nil
}

//关闭
func (s *TcpProxy) Close() {
	s.Stop()
}

//暂停接受新连接
func (s *TcpProxy) Stop() {
	if s.closeChan == nil {
		return
	}

	select {
	case <-s.closeChan:
		//已经关闭了
		return
	default:
		close(s.closeChan)
	}

	logger.Info("转发代理服停止")
	_ = s.listener.Close()
}

//TCP服务器主循环
func (s *TcpProxy) Run() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(r)
				logger.Error("TCP服务器主循环退出")
			}
		}()

		for {
			err := s.accept()
			if err != nil {
				select {
				case <-s.closeChan:
					//已经关闭了
					return
				}
			}
		}
	}()
}

//接收一个连接 协程里
func (s *TcpProxy) accept() error {
	conn, err := s.listener.Accept()
	if err != nil {
		return err
	}

	go s.beginServe(conn)

	return nil
}

//创建连接
func (s *TcpProxy) beginServe(src net.Conn) {
	defer src.Close()

	logger.Info("转发代理服 建立连接 ", src.RemoteAddr())
	dst, err := net.Dial("tcp", s.addrTarget)
	if err != nil {
		return
	}

	exit := make(chan bool, 1)
	go func(src net.Conn, dst net.Conn, exit chan bool) {
		io.Copy(dst, src)
		exit <- true
	}(src, dst, exit)
	go func(src net.Conn, dst net.Conn, exit chan bool) {
		io.Copy(src, dst)
		exit <- true
	}(src, dst, exit)
	<-exit
	dst.Close()
}
