package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
	"silly/logger"
	"silly/transport"
	"strings"
)

//for register of remote function
type RemoteFunReg struct {
	objTable map[interface{}]map[string][]reflect.Type
}

func (r *RemoteFunReg) Init() {
	r.objTable = make(map[interface{}]map[string][]reflect.Type)
}

func (r *RemoteFunReg) RegisterRemoteFunc(obj interface{}) {
	objType := reflect.TypeOf(obj)
	tag := "RemoteCall"
	for i := 0; i < objType.NumMethod(); i++ {
		method := objType.Method(i)
		if !strings.HasPrefix(method.Name, tag) {
			continue
		}
		paramNum := method.Type.NumIn()
		if paramNum <= 0 {
			continue
		}
		paramType := make([]reflect.Type, paramNum-1)
		for j := 1; j < paramNum; j++ {
			paramType[j-1] = method.Type.In(j)
		}
		if r.objTable[objType] == nil {
			r.objTable[objType] = make(map[string][]reflect.Type)
		}
		if r.objTable[objType][method.Name] == nil {
			r.objTable[objType][method.Name] = paramType
		}
	}
}

func (r *RemoteFunReg) GetMethodParams(obj interface{}, funName string) []reflect.Type {
	objType := reflect.TypeOf(obj)
	if r.objTable == nil || r.objTable[objType] == nil || r.objTable[objType][funName] == nil {
		return nil
	}
	return r.objTable[objType][funName]
}

func ReflectCallMethod(obj interface{}, methodName string, params ...interface{}) []reflect.Value {
	valueObj := reflect.ValueOf(obj)
	method := valueObj.MethodByName(methodName)
	if !method.IsValid() {
		logger.Info("反射调用函数，函数不存在： ", methodName)
		return nil
	}
	param := make([]reflect.Value, len(params))
	for i, v := range params {
		param[i] = reflect.ValueOf(v)
	}
	return method.Call(param)
}

func UnmarshalServerCallMsg(data []byte, reg *RemoteFunReg, obj interface{}) (n string, p []interface{}, e error) {
	buf := bytes.NewBuffer(data)
	msgId := ReadInt32(buf)
	logger.Info("msgId: ", msgId)
	methodName := ReadString(buf)
	paramTypes := reg.GetMethodParams(obj, methodName)
	if paramTypes == nil {
		e = errors.New(fmt.Sprintf("未注册的函数：%v", methodName))
		return
	}
	params := make([]interface{}, len(paramTypes))
	for i, t := range paramTypes {
		param := reflect.New(t).Interface()
		paramBuf := ReadBytes(buf)
		err := msgpack.Unmarshal(paramBuf, param)
		if err != nil {
			e = errors.New(fmt.Sprintf("远程调用参数错误,函数: %v, err: %v", methodName, err))
			return
		}
		params[i] = reflect.ValueOf(param).Elem().Interface()
	}
	return methodName, params, nil
}

func WriteUint16(buf *bytes.Buffer, n uint16) {
	b := make([]byte, 2)
	transport.DefaultEndian.PutUint16(b, n)
	buf.Write(b)
}

func ReadUint16(buf *bytes.Buffer) uint16 {
	if buf.Len() < 2 {
		panic("ReadUint16, buff is too short")
	}
	n := transport.DefaultEndian.Uint16(buf.Next(2))
	return n
}

func WriteInt32(buf *bytes.Buffer, n int) {
	b := make([]byte, 4)
	transport.DefaultEndian.PutUint32(b, uint32(n))
	buf.Write(b)
}

func ReadInt32(buf *bytes.Buffer) int {
	if buf.Len() < 4 {
		panic("ReadInt, buff is too short")
	}
	n := int(transport.DefaultEndian.Uint32(buf.Next(4)))
	return n
}

func ModifyInt32(buf *bytes.Buffer, idx int, n int) {
	if buf.Len() < idx+4 {
		panic("ModifyInt32, buff is too short")
	}
	b := make([]byte, 4)
	transport.DefaultEndian.PutUint32(b, uint32(n))
	copy(buf.Bytes()[idx:], b)
}

func WriteString(buf *bytes.Buffer, s string) {
	b := []byte(s)
	WriteInt32(buf, len(b))
	buf.Write(b)
}

func ReadString(buf *bytes.Buffer) string {
	len := ReadInt32(buf)
	if buf.Len() < len {
		panic("ReadString, buff is too short")
	}
	s := string(buf.Next(len))
	return s
}

func WriteBytes(buf *bytes.Buffer, b []byte) {
	WriteInt32(buf, len(b))
	buf.Write(b)
}

func ReadBytes(buf *bytes.Buffer) []byte {
	len := ReadInt32(buf)
	if buf.Len() < len {
		panic("ReadBytes, buff is too short")
	}
	b := buf.Next(len)
	return b
}
