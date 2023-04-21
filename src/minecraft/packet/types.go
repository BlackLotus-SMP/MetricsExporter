package packet

import (
	"errors"
	"io"
)

type Field interface {
	FieldEncoder
	FieldDecoder
}

type FieldEncoder io.WriterTo

type FieldDecoder io.ReaderFrom

type (
	Byte          int8
	UnsignedShort uint16
	String        string
	VarInt        int32
)

const MaxVarIntLen = 5

func (s String) WriteTo(w io.Writer) (int64, error) {
	byteStr := []byte(s)
	n1, err := VarInt(len(byteStr)).WriteTo(w)
	if err != nil {
		return n1, err
	}
	n2, err := w.Write(byteStr)
	return n1 + int64(n2), err
}

func (b Byte) WriteTo(w io.Writer) (n int64, err error) {
	nn, err := w.Write([]byte{byte(b)})
	return int64(nn), err
}

func (us UnsignedShort) WriteTo(w io.Writer) (int64, error) {
	n := uint16(us)
	nn, err := w.Write([]byte{byte(n >> 8), byte(n)})
	return int64(nn), err
}

func (v VarInt) WriteTo(w io.Writer) (n int64, err error) {
	var vi [MaxVarIntLen]byte
	nn := v.WriteToBytes(vi[:])
	nn, err = w.Write(vi[:nn])
	return int64(nn), err
}

func readByte(r io.Reader) (int64, byte, error) {
	if r, ok := r.(io.ByteReader); ok {
		v, err := r.ReadByte()
		return 1, v, err
	}
	var v [1]byte
	n, err := r.Read(v[:])
	return int64(n), v[0], err
}

// WriteToBytes encodes the VarInt into buf and returns the number of bytes written.
// If the buffer is too small, WriteToBytes will panic.
func (v VarInt) WriteToBytes(buf []byte) int {
	num := uint32(v)
	i := 0
	for {
		b := num & 0x7F
		num >>= 7
		if num != 0 {
			b |= 0x80
		}
		buf[i] = byte(b)
		i++
		if num == 0 {
			break
		}
	}
	return i
}

func (v *VarInt) ReadFrom(r io.Reader) (n int64, err error) {
	var V uint32
	var num, n2 int64
	for sec := byte(0x80); sec&0x80 != 0; num++ {
		if num > MaxVarIntLen {
			return n, errors.New("VarInt is too big")
		}

		n2, sec, err = readByte(r)
		n += n2
		if err != nil {
			return n, err
		}

		V |= uint32(sec&0x7F) << uint32(7*num)
	}

	*v = VarInt(V)
	return
}

func (v VarInt) Len() int {
	switch {
	case v < 0:
		return MaxVarIntLen
	case v < 1<<(7*1):
		return 1
	case v < 1<<(7*2):
		return 2
	case v < 1<<(7*3):
		return 3
	case v < 1<<(7*4):
		return 4
	default:
		return 5
	}
}
