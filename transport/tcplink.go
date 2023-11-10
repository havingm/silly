package transport

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"silly/logger"
	"sync"
	"time"
)

const (
	WriteBuffLen  = 500
	TimeoutSecond = 5
)

type TcpLink struct {
	sync.RWMutex
	conn         net.Conn
	linkId       uint64
	serviceId    uint
	buffer       chan *NetPack
	readDeadline time.Duration
	sendTimeout  time.Duration
	pingTime     time.Duration
	service      ITcpSocket
}

type ITcpLink interface {
	GetLinkId() uint64
	Send([]byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	RemoteIp() string
	Close()
	Stop()
}

func NewTcpLink(conn net.Conn, service ITcpSocket, ping int) *TcpLink {
	l := &TcpLink{}
	l.conn = conn
	l.service = service
	l.buffer = make(chan *NetPack, WriteBuffLen)
	l.sendTimeout = TimeoutSecond * time.Second
	l.pingTime = time.Duration(ping) * time.Second
	return l
}

func (l *TcpLink) Run() {
	if l.conn == nil {
		logger.Error("TcpLink.Run, conn is nil...")
		return
	}
	//数据写入
	go func() {
		defer func() {
			logger.Info("数据写入线程退出")
			if r := recover(); r != nil {
				logger.Error(r)
				l.Close()
			}
		}()

		l.RLock()
		conn := l.conn
		l.RUnlock()

		for conn != nil {
			pack, open := <-l.buffer
			if !open {
				return
			}
			if pack.Stop {
				l.Close()
				return
			}
			//conn.SetWriteDeadline(time.Now().Add(20 * time.Second))
			_, err := conn.Write(pack.ConvertBytes(4))
			if err != nil {
				logger.Error(err)
				l.Close()
				return
			}
		}
	}()
	//数据接收
	go func() {
		defer func() {
			logger.Info("数据接收线程退出")
			if r := recover(); r != nil {
				logger.Error(r)
				l.Close()
			}
		}()

		l.RLock()
		conn := l.conn
		l.RUnlock()

		readErr := 0
		for conn != nil {
			lenBuf := make([]byte, 4)
			if l.readDeadline > 0 {
				l.conn.SetReadDeadline(time.Now().Add(l.readDeadline))
			}
			if _, err := io.ReadFull(conn, lenBuf); err != nil {
				//EOF表示连接已关闭
				if err.Error() != "EOF" {
					logger.Warning(err)
				}
				l.Close()
				return
			}
			dataLen := binary.LittleEndian.Uint32(lenBuf)
			//ping 包
			if dataLen == 0 {
				logger.Info("收到ping包, linkId: ", l.linkId)
				continue
			}
			//读取数据
			data := make([]byte, dataLen)
			_, err := io.ReadFull(conn, data)
			if err != nil {
				if readErr >= 100 {
					l.Close()
					return
				}
				readErr++
				continue
			}
			l.service.OnLinkRecv(l, data)
		}
	}()
	//ping
	if l.pingTime > 0 {
		pingBuf := make([]byte, 4)
		binary.LittleEndian.PutUint32(pingBuf, 0)
		go func() {
			defer func() {
				logger.Info("ping线程退出")
				if r := recover(); r != nil {
					logger.Error(r)
					l.Close()
				}
			}()

			l.RLock()
			conn := l.conn
			l.RUnlock()

			for {
				time.Sleep(l.pingTime)
				_, err := conn.Write(pingBuf)
				if err != nil {
					l.Close()
					return
				}
			}
		}()
	}
}

func (l *TcpLink) GetLinkId() uint64 {
	return l.linkId
}

func (l *TcpLink) Send(data []byte) error {
	if l.conn == nil {
		return errors.New("TcpLink Send, conn is nil")
	}
	if data == nil {
		return errors.New("TcpLink Send, data is nil")
	}
	timeout := time.NewTimer(l.sendTimeout)
	defer timeout.Stop()
	pack := &NetPack{
		Len:  uint32(len(data)),
		Data: data,
	}
	select {
	//todo 缓冲区满了会阻塞
	case l.buffer <- pack:
	case <-timeout.C:
		l.Close()
		logger.Warning("TcpLink Send timeout, closed")
	}
	return nil
}

func (l *TcpLink) LocalAddr() net.Addr {
	return nil
}

func (l *TcpLink) RemoteAddr() net.Addr {
	return nil
}

func (l *TcpLink) RemoteIp() string {
	return ""
}

func (l *TcpLink) Close() {
	l.Lock()
	defer l.Unlock()
	if l.conn == nil {
		return
	}
	conn := l.conn.(*net.TCPConn)
	l.conn = nil
	close(l.buffer)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(r)
			}
		}()
		l.service.OnLinkClose(l)
		//优雅关闭时间
		conn.SetLinger(5)
		conn.Close()
	}()
}

func (l *TcpLink) Stop() {
	timer := time.NewTimer(100 * time.Millisecond)
	defer timer.Stop()
	stopPack := &NetPack{Stop: true}
	select {
	case l.buffer <- stopPack:
	case <-timer.C:
		logger.Warning("TcpLink.Stop timeout")
		l.Close()
	}
}
