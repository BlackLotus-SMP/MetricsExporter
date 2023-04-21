package minecraft

import (
	"context"
	"fmt"
	"io"
	"metrics-exporter/src/minecraft/packet"
	"net"
)

type MCDialer interface {
	DialMCContext(ctx context.Context, mcAddr string, port uint) (*WrappedConn, error)
}

type Dialer net.Dialer

type WrappedConn struct {
	Socket net.Conn
	io.Reader
	io.Writer
}

func (d *Dialer) DialMCContext(ctx context.Context, mcAddr string, port uint) (*WrappedConn, error) {
	conn, err := (*net.Dialer)(d).DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", mcAddr, port))
	if err != nil {
		return nil, err
	}
	return &WrappedConn{
		Socket: conn,
		Reader: conn,
		Writer: conn,
	}, nil
}

func (c *WrappedConn) WritePacket(p packet.Packet) error {
	return p.Pack(c.Writer)
}

func (c *WrappedConn) ReadPacket(p *packet.Packet) error {
	return p.UnPack(c.Reader)
}
