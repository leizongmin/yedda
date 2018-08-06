package protocol

import (
	"encoding/binary"
)

func Uint32ToBytes(v uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	return b
}

func BytesToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}
