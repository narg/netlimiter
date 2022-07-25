package netlimiter

import (
	"context"
	"io/ioutil"
	"log"
	"net"

	// "reflect"
	"testing"
)

func TestConn(t *testing.T) {
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

	var globalLimit, connLimit int64
	globalLimit = 15
	connLimit = 5

	ln := NewListener(l, ctx, globalLimit, connLimit)

	conn, err := ln.Accept()
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 10)

	_, err = conn.Read(buf)
	if err != nil {
		log.Println(err)
		// t.Fatal(err)
	}

	_, err = conn.Write([]byte("Test Response\n"))
	if err != nil {
		log.Println(err)
	}

	conn.SetConfig(globalLimit - 1)
	if conn.Config().Conn() != globalLimit-1 {
		t.Error("Conn limit is not equal to expected value")
	}

	conn.SetConfig(globalLimit + 1)
	if conn.Config().Conn() != globalLimit {
		t.Error("Conn limit is not equal to expected value")
	}

	defer conn.Close()

}
