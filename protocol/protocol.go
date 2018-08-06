package protocol

import (
	"encoding/binary"
	"io"
)

type Package struct {
	Version uint16 `json:"version"` // 版本
	Op      OpType `json:"op"`      // 操作类型
	Length  uint16 `json:"length"`  // 数据长度
	Data    []byte `json:"data"`    // 数据内容
}

type OpType uint8

const (
	_ OpType = iota
	OpPing
	OpPong
	OpGet
	OpGetResult
	OpIncr
	OpIncrResult
)

const CurrentVersion = 1

func NewPackage(version uint16, op OpType, data []byte) *Package {
	return &Package{
		Version: version,
		Op:      op,
		Length:  uint16(len(data)),
		Data:    data,
	}
}

func NewPackageFromReader(r io.Reader) (p *Package, err error) {
	p = &Package{}
	err = p.UnPack(r)
	return p, err
}

func PackToWriter(w io.Writer, version uint16, op OpType, data []byte) error {
	return NewPackage(version, op, data).Pack(w)
}

func (p *Package) Pack(w io.Writer) (err error) {
	b := make([]byte, 5)
	binary.BigEndian.PutUint16(b, p.Version)
	b[2] = byte(p.Op)
	binary.BigEndian.PutUint16(b[3:], p.Length)
	w.Write(b)
	w.Write(p.Data)
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
