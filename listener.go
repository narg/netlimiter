package netlimiter

import (
	"context"
	"log"
	"net"
)

type Listener struct {
	net.Listener
	ctx    context.Context
	config *Config
}

// Return *Config
func (ln *Listener) Config() *Config {
	return ln.config
}

// Set global and connection bandwidth limit
func (ln *Listener) SetConfig(global int64, conn int64) {
	ln.config.SetGlobal(global)
	ln.config.SetGlobalLimiter()

	ln.config.SetConn(conn)
	ln.config.SetConnLimiter()
	log.Printf("Listener Config has been updated!\tGlobal: %d,\tConn: %d", ln.config.Global(), ln.config.Conn())
}

// Accept request and return new connection
func (ln *Listener) Accept() (*Conn, error) {
	conn, err := ln.Listener.Accept()
	if err != nil {
		var c *Conn
		return c, err
	}

	return NewConn(conn, ln.ctx, ln.config)
}

func (ln *Listener) Close() error {
	return ln.Listener.Close()
}

// Create and return new listener
func NewListener(listener net.Listener, ctx context.Context, global int64, conn int64) *Listener {
	conf := NewConfig(global, conn)

	var ln Listener

	ln.Listener = listener
	ln.ctx = ctx
	ln.config = conf

	return &ln
}
