package wsclient

import "sync/atomic"

var maxWriterBuffer int64 = 4

func SetMaxWriterBuffer(num int64) {
	atomic.StoreInt64(&maxWriterBuffer, num)
}

func GetMaxWriterBuffer() int64 {
	return atomic.LoadInt64(&maxWriterBuffer)
}
