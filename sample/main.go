package main

import (
	"bufio"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
	"os"
	"silly/logger"
	"silly/msg"
	"silly/transport"
	"strconv"
	"strings"
	"time"
)

type ClientManager struct {
}

func (c *ClientManager) OnLinkRecv(link *transport.TcpLink, data []byte) {
	t2 := new(msg.T2)
	err := proto.Unmarshal(data, t2)
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info("client recv: ", t2)
	}
}

func (c *ClientManager) OnLinkClose(link *transport.TcpLink) {
	logger.Info("client closed")
}

var client *transport.TcpLink

func main() {

	//urlValues := url.Values{
	//	"s2s":         {"1"},
	//	"app_token":   {"s2qlwq6vxgcg"},
	//	"event_token": {"5380so"},
	//	"adid":        {"68010b3102a2aaafe6e0d4e6236689d0"},
	//	"sdk_open_id": {"100001"},
	//	"username":    {"having"},
	//}
	//urlValues["draw_type"] = []string{"0"}
	//urlValues["draw_cnt"] = []string{"10"}
	//resp, err := http.PostForm("https://s2s.adjust.com/event", urlValues)
	//if err != nil {
	//	logger.Error("postForm, err: ", err)
	//} else {
	//
	//	body, err := ioutil.ReadAll(resp.Body)
	//	if err != nil {
	//		logger.Error(err)
	//	} else {
	//		result := string(body)
	//		var data map[string]interface{}
	//		logger.Info("result: ", result)
	//		if err := json.Unmarshal([]byte(result), &data); err != nil {
	//			logger.Error("err: ", err)
	//		} else {
	//			var status = data["status"].(string)
	//			if "OK" != status {
	//				logger.Error(" OK!= status")
	//			}
	//		}
	//	}
	//
	//}

	logger.Info("Test started...")
	service := transport.NewTcpService("Gate")
	err := service.Start("localhost:8080", 0)
	if err != nil {
		logger.Warning(err)
		reader := bufio.NewReader(os.Stdin)
		reader.ReadLine()
		return
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	mgr := &ClientManager{}

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		cmd := string(data)
		strings.Split(cmd, " ")
		switch cmd {
		case "quit":
			logger.Info("quited........")
			return
		default:
			logger.Info(fmt.Sprintf("收到命令：[%v]", cmd))
			cmdSplit := strings.Split(cmd, " ")
			splitLen := len(cmdSplit)
			if splitLen <= 0 {
				logger.Info("can't split cmd...")
				continue
			}
			switch cmdSplit[0] {
			case "make":
				num := 1
				if splitLen > 1 {
					num, _ = strconv.Atoi(cmdSplit[1])
				}
				for i := 0; i < num; i++ {
					conn, err := net.DialTCP("tcp", nil, tcpAddr)
					if err != nil {
						logger.Error("DialTcp failed cause ", err)
					}
					client = transport.NewTcpLink(conn, mgr, 10)
					client.Run()
					time.Sleep(10)
				}
			case "send":
				if splitLen <= 1 {
					logger.Warning("send nothing...")
					continue
				}
				client.Send([]byte(cmdSplit[1]))
			}
		case "close":
			if client != nil {
				client.Stop()
				client = nil
			}
		case "lsend":
			idx := 1
			for client != nil {
				client.Send([]byte(fmt.Sprintf("loop send: %d", idx)))
				idx++
				time.Sleep(10)
			}
		case "psend":
			m1 := &msg.T1{}
			m1.B = proto.Bool(true)
			m1.I = proto.Int32(100)
			m1.S = proto.String("Test")
			m2 := &msg.T2{}
			m2.R = make([]*msg.T1, 0)
			for i := 1; i <= 10; i++ {
				t := &msg.T1{I: proto.Int32(int32(i)),
					B: proto.Bool(i%2 == 0),
					S: proto.String(strconv.Itoa(i))}
				m2.R = append(m2.R, t)
			}
			data, _ := proto.Marshal(m2)
			if client == nil {
				logger.Error("client is nil...")
				break
			}
			client.Send(data)

		}
	}

}
