package netlimiter

import (
	"context"
	"io/ioutil"
	"log"
	"net"

	// "reflect"
	"testing"
)

func TestListen(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	go func() {
		conn, err := net.Dial("tcp", ":3000")
		if err != nil {
			log.Println(err)
		}
		defer conn.Close()
	}()

	ctx := context.Background()
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	var global, conn int64
	global = 15
	conn = 5

	ln := NewListener(l, ctx, global, conn)

	if ln.Config().Global() != global {
		t.Error("Global limit is not equal to expected value")
	}

	if ln.Config().Conn() != conn {
		t.Error("Conn limit is not equal to expected value")
	}

	global = 34
	conn = 7
	ln.SetConfig(global, conn)
	if ln.Config().Global() != global {
		t.Error("Global limit is not equal to expected value")
	}

	if ln.Config().Conn() != conn {
		t.Error("Conn limit is not equal to expected value")
	}
}
