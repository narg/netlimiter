package netlimiter

import "math"
import "testing"

func TestNewConfig(t *testing.T) {
	var global, conn int64

	conn = int64(100)
	global = int64(200)

	config := NewConfig(global, conn)

	if config.Global() != global {
		t.Error("got: ", config.Global(), "\texpected: ", global)
	}

	if config.Conn() != conn {
		t.Error("got: ", config.Conn(), "\texpected: ", conn)
	}
}

func TestSetGlobal(t *testing.T) {
	var global, conn int64

	conn = int64(100)
	global = int64(200)

	config := NewConfig(global, conn)

	var new int64

	//
	new = 400
	config.SetGlobal(new)

	if config.Global() != new {
		t.Error("got: ", config.Global(), "\texpected: ", new)
	}

	//
	new = math.MaxInt64
	config.SetGlobal(new)

	if config.Global() != new {
		t.Error("got: ", config.Global(), "\texpected: ", new)
	}

	//
	new = 0
	config.SetGlobal(new)

	if config.Global() != new {
		t.Error("got: ", config.Global(), "\texpected: ", math.MaxInt64)
	}
}

func TestSetConn(t *testing.T) {
	var global, conn int64

	conn = int64(100)
	global = int64(200)

	config := NewConfig(global, conn)

	var new int64

	//
	new = config.Global() - 1
	config.SetConn(new)

	if config.Conn() != new {
		t.Error("got: ", config.Conn(), "\texpected: ", new)
	}

	//
	new = config.Global() + 1
	config.SetConn(new)

	if config.Conn() != config.Global() {
		t.Error("got: ", config.Conn(), "\texpected: ", config.Global())
	}

	//
	new = 0
	config.SetConn(new)

	if config.Conn() != config.Global() {
		t.Error("got: ", config.Conn(), "\texpected: ", config.Global())
	}
}
