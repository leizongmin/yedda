package core

// MD5键类型
type HashKey [16]byte

// Map: MD5键->原始字符
type HashKeyMap map[HashKey][]byte
