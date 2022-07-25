package netlimiter

import (
	"math"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Config struct {
	global        int64
	conn          int64
	globalLimiter *rate.Limiter
	connLimiter   *rate.Limiter
	locker        sync.Mutex
}

// Return global bandwidth limit
func (c *Config) Global() int64 {
	c.locker.Lock()
	defer c.locker.Unlock()

	return c.global
}

// Set global bandwidth limit
func (c *Config) SetGlobal(n int64) {
	c.locker.Lock()
	defer c.locker.Unlock()

	if n < 1 {
		c.global = math.MaxInt64
	}

	c.global = n
}

// Return connection bandwidth limit
func (c *Config) Conn() int64 {
	c.locker.Lock()
	defer c.locker.Unlock()

	return c.conn
}

// Set connection bandwidth limit
func (c *Config) SetConn(n int64) {
	c.locker.Lock()
	defer c.locker.Unlock()

	if n > 0 && n < c.global {
		c.conn = n
	} else {
		c.conn = c.global
	}
}

// Return global bandwith limiter
func (c *Config) GlobalLimiter() *rate.Limiter {
	return c.globalLimiter
}

// Set global bandwith limiter
func (c *Config) SetGlobalLimiter() {
	c.locker.Lock()
	defer c.locker.Unlock()

	limit := float64(c.global)
	// max size limit
	burst := 5 * 60 * int(c.global)

	c.globalLimiter = rate.NewLimiter(rate.Every(time.Second/time.Duration(limit)), burst)
}

// Return connection bandwith limiter
func (c *Config) ConnLimiter() *rate.Limiter {
	return c.connLimiter
}

// Set connection bandwith limiter
func (c *Config) SetConnLimiter() {
	c.locker.Lock()
	defer c.locker.Unlock()

	limit := float64(c.conn)
	// max size limit
	burst := 5 * 60 * int(c.global)

	c.connLimiter = rate.NewLimiter(rate.Every(time.Second/time.Duration(limit)), burst)
}

// Create and return new config
func NewConfig(global int64, conn int64) *Config {
	var conf Config

	conf.SetGlobal(global)
	conf.SetConn(conn)

	conf.SetGlobalLimiter()
	conf.SetConnLimiter()

	return &conf
}
