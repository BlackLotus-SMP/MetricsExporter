package packet

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

type Packet struct {
	ID   int32
	Data []byte
}

func Marshal[ID ~int32 | int](id ID, fields ...FieldEncoder) (pk Packet) {
	var pb Builder
	for _, v := range fields {
		pb.WriteField(v)
	}
	return pb.Packet(int32(id))
}

var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func (p *Packet) Pack(w io.Writer) error {
	buffer := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buffer)
	buffer.Reset()

	Length := VarInt(VarInt(p.ID).Len() + len(p.Data))
	_, _ = Length.WriteTo(buffer)

	_, _ = VarInt(p.ID).WriteTo(buffer)
	buffer.Write(p.Data)

	_, err := w.Write(buffer.Bytes())
	return err
}

func (p *Packet) UnPack(r io.Reader) error {
	var Length VarInt
	_, err := Length.ReadFrom(r)
	if err != nil {
		return err
	}

	var PacketID VarInt
	n, err := PacketID.ReadFrom(r)
	if err != nil {
		return err
	}
	p.ID = int32(PacketID)

	lengthOfData := int(Length) - int(n)
	if lengthOfData < 0 {
		return fmt.Errorf("uncompressed packet error: length is %d", lengthOfData)
	}
	if cap(p.Data) < lengthOfData {
		p.Data = make([]byte, lengthOfData)
	} else {
		p.Data = p.Data[:lengthOfData]
	}
	_, err = io.ReadFull(r, p.Data)
	if err != nil {
		return err
	}
	return nil
}
