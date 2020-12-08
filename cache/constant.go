package cache

import "unsafe"

const (
	LOCK_REDIS = "incr_redis"
)

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
