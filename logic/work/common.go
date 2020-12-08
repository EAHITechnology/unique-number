package work

import "unsafe"

func workInfo2Byte(info workInfo) []byte {
	len := unsafe.Sizeof(info)
	byt := &tempAddr{
		addr: uintptr(unsafe.Pointer(&info)),
		cap:  int(len),
		len:  int(len),
	}
	return *(*[]byte)(unsafe.Pointer(byt))
}

func byte2WorkInfo(res []byte) *workInfo {
	return *(**workInfo)(unsafe.Pointer(&res))
}
