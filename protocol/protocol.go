package protocol

import (
	"encoding/binary"
	"io"
)

type Package struct {
	Version uint16 `json:"version"` // 版本
	Op      uint8  `json:"op"`      // 操作类型
	Length  uint16 `json:"length"`  // 数据长度
	Data    []byte `json:"data"`    // 数据内容
}

type OpType uint8

const (
	_ OpType = iota
	OpPing
	OpPong
	OpGet
	OpIncr
)

func NewPackage(version uint16, op OpType, data []byte) *Package {
	return &Package{
		Version: version,
		Op:      uint8(op),
		Length:  uint16(len(data)),
		Data:    data,
	}
}

func NewPackageFromReader(r io.Reader) (p *Package, err error) {
	p = &Package{}
	err = p.UnPack(r)
	return p, err
}

func (p *Package) Pack(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, &p.Version)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, &p.Op)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, &p.Length)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, &p.Data)
	return err
}

func (p *Package) UnPack(r io.Reader) (err error) {
	err = binary.Read(r, binary.BigEndian, &p.Version)
	if err != nil {
		return err
	}
	err = binary.Read(r, binary.BigEndian, &p.Op)
	if err != nil {
		return err
	}
	err = binary.Read(r, binary.BigEndian, &p.Length)
	if err != nil {
		return err
	}
	p.Data = make([]byte, p.Length)
	err = binary.Read(r, binary.BigEndian, &p.Data)
	return err
}
