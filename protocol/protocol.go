package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Package struct {
	Version uint16 `json:"version"` // 版本
	Op      OpType `json:"op"`      // 操作类型
	Length  uint16 `json:"length"`  // 数据长度
	Data    []byte `json:"data"`    // 数据内容
}

type OpType uint16

const (
	OpPing       = 0x1
	OpPong       = 0x2
	OpGet        = 0x3
	OpGetResult  = 0x4
	OpIncr       = 0x5
	OpIncrResult = 0x6
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
	b := make([]byte, 6)
	binary.BigEndian.PutUint16(b, p.Version)
	binary.BigEndian.PutUint16(b[2:], uint16(p.Op))
	binary.BigEndian.PutUint16(b[4:], p.Length)
	w.Write(b)
	w.Write(p.Data)
	return err
}

func (p *Package) UnPack(r io.Reader) (err error) {
	b := make([]byte, 6)
	n, err := r.Read(b)
	if err != nil {
		return err
	}
	if n != 6 {
		return fmt.Errorf("expected to read %d bytes but got %d bytes", 6, n)
	}
	p.Version = binary.BigEndian.Uint16(b)
	p.Op = OpType(binary.BigEndian.Uint16(b[2:]))
	p.Length = binary.BigEndian.Uint16(b[4:])
	p.Data = make([]byte, p.Length)
	if p.Length > 0 {
		n, err = r.Read(p.Data)
		if n != 6 {
			return fmt.Errorf("expected to read %d bytes but got %d bytes", 6, p.Length)
		}
	}
	return err
}
