package transport

import (
	"encoding/binary"
	"silly/logger"
)

type NetPack struct {
	Len  uint32
	Data []byte
	Stop bool
}

func (p *NetPack) ConvertBytes(packLen int) (buff []byte) {
	if packLen == 2 {
		if p.Len != uint32(uint16(p.Len)) {
			logger.Error("数据长度超出限制: ", p.Len)
		}
		buff = make([]byte, 2+uint16(len(p.Data)))
		binary.LittleEndian.PutUint16(buff, uint16(p.Len))
		copy(buff[2:], p.Data)
	} else if packLen == 4 {
		buff = make([]byte, 4+uint32(len(p.Data)))
		binary.LittleEndian.PutUint32(buff, p.Len)
		copy(buff[4:], p.Data)
	} else {
		panic("不支持的包长度")
	}
	return buff
}
