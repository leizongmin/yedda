package service

import (
	"bytes"
	"encoding/binary"
)

// 调用参数
type CmdArg struct {
	// 数据库号，>= 0
	Db uint32
	// 命名空间
	Ns string
	// 有效期
	Milliseconds uint32
	// 键名
	Key []byte
	// 增加数量
	Count uint32
}

// 创建命令调用参数对象
func NewCmdArg(db uint32, ns string, milliseconds uint32, key []byte, count uint32) *CmdArg {
	return &CmdArg{
		Db:           db,
		Ns:           ns,
		Milliseconds: milliseconds,
		Key:          key,
		Count:        count,
	}
}

// 生成 []byte 数据
func (a *CmdArg) Bytes() ([]byte, error) {
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	ns := []byte(a.Ns)
	if len(ns) >= 256 {
		ns = ns[0:256]
	}
	if len(a.Key) >= 256 {
		a.Key = a.Key[0:256]
	}

	nsLen := uint8(len(ns))
	keyLen := uint8(len(a.Key))

	err = binary.Write(buf, binary.BigEndian, &a.Db)
	err = binary.Write(buf, binary.BigEndian, &nsLen)
	err = binary.Write(buf, binary.BigEndian, &ns)
	err = binary.Write(buf, binary.BigEndian, &a.Milliseconds)
	err = binary.Write(buf, binary.BigEndian, &keyLen)
	err = binary.Write(buf, binary.BigEndian, &a.Key)
	err = binary.Write(buf, binary.BigEndian, &a.Count)
	return buf.Bytes(), err
}

// 通过 []byte 数据生成参数结构
func NewCmdArgFromBytes(b []byte) (*CmdArg, error) {
	var err error
	buf := bytes.NewBuffer(b)
	a := CmdArg{}
	var nsLen, keyLen uint8
	err = binary.Read(buf, binary.BigEndian, &a.Db)
	err = binary.Read(buf, binary.BigEndian, &nsLen)
	ns := make([]byte, nsLen)
	err = binary.Read(buf, binary.BigEndian, &ns)
	a.Ns = string(ns)
	err = binary.Read(buf, binary.BigEndian, &a.Milliseconds)
	err = binary.Read(buf, binary.BigEndian, &keyLen)
	a.Key = make([]byte, keyLen)
	err = binary.Read(buf, binary.BigEndian, &a.Key)
	err = binary.Read(buf, binary.BigEndian, &a.Count)
	return &a, err
}
