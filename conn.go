package netlimiter

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Conn struct {
	net.Conn
	ctx    context.Context
	config *Config
}

// Run *rate.Limiter Wait for Token consumptions
func waitN(limiter *rate.Limiter, ctx context.Context, n int, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Printf("Limiter Values\tLimit: %f\tBurst: %d\tn: %d", limiter.Limit(), limiter.Burst(), n)

	err := limiter.WaitN(ctx, n)
	if err != nil {
		return err
	}

	return nil
}

// Read reads data from the connection.
func (c *Conn) Read(n []byte) (int, error) {
	start := time.Now()
	var wg sync.WaitGroup

	wg.Add(2)
	go waitN(c.config.globalLimiter, c.ctx, len(n), &wg)
	go waitN(c.config.connLimiter, c.ctx, len(n), &wg)

	wg.Wait()

	elapsed := time.Since(start)

	log.Printf("Read Elapsed Time: %s", elapsed)

	return c.Conn.Read(n)
}

// Write writes data to the connection.
func (c *Conn) Write(n []byte) (int, error) {
	start := time.Now()
	var wg sync.WaitGroup

	wg.Add(2)
	go waitN(c.config.globalLimiter, c.ctx, len(n), &wg)
	go waitN(c.config.connLimiter, c.ctx, len(n), &wg)

	wg.Wait()

	elapsed := time.Since(start)

	log.Printf("Write Elapsed Time: %s", elapsed)

	return c.Conn.Write(n)
}

// Close closes the connection.
func (c *Conn) Close() error {
	return c.Conn.Close()
}

// Return *Config
func (c *Conn) Config() *Config {
	return c.config
}

// Set connection bandwidth limit
func (c *Conn) SetConfig(conn int64) {
	c.config.SetConn(conn)
	c.config.SetConnLimiter()
	log.Printf("Conn Config has been updated!\tGlobal: %d,\tConn: %d", c.config.Global(), c.config.Conn())
}

// Create and return new connection
func NewConn(conn net.Conn, ctx context.Context, conf *Config) (*Conn, error) {
	var c Conn

	c.Conn = conn
	c.ctx = ctx
	c.config = conf

	return &c, nil
}
