package silly

import (
	"bytes"
	"github.com/vmihailenco/msgpack/v5"
	. "silly/utils"
)

var msgId uint32

func EncodeServerMsg(methodName string, params ...interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	//消息长度
	msgId++
	WriteInt32(buf, int(msgId))
	WriteString(buf, methodName)
	for _, param := range params {
		pb, err := msgpack.Marshal(param)
		if err != nil {
			return nil, err
		}
		WriteBytes(buf, pb)
	}
	return buf.Bytes(), nil
}
