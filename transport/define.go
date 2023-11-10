package transport

import "sync/atomic"

var _increaseTcpLinkId_ uint64

func AutoTcpLinkId() uint64 {
	return atomic.AddUint64(&_increaseTcpLinkId_, 1)
}
