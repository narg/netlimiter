# Net Limiter

Limit TCP Read & Write bandwith


## Example

```go
package main

import (
	"context"
	"log"
	"net"

	"github.com/narg/netlimiter"
)

func main() {
	ctx := context.Background()
	tcpListener, err := net.Listen("tcp", ":5030")
	if err != nil {
		panic(err)
	}

	var globalLimit, connLimit int64

	globalLimit = 15
	connLimit = 5

	listener := netlimiter.NewListener(tcpListener, ctx, globalLimit, connLimit)

	host, port, err := net.SplitHostPort(tcpListener.Addr().String())
	if err != nil {
		panic(err)
	}

	log.Printf("Listening on host: %s, port: %s\n", host, port)

	requestNumber := 0
	readBufSize := 10

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		requestNumber++

		go func(conn *netlimiter.Conn) {
			defer conn.Close()

			log.Printf("Request #%d", requestNumber)

			buf := make([]byte, readBufSize)

			_, err := conn.Read(buf)
			if err != nil {
				log.Println(err)
			}

			_, err = conn.Write([]byte("Test Response\n"))
			if err != nil {
				log.Println(err)
			}

			if requestNumber%5 == 0 {
				connLimit *= 2
				conn.SetConfig(connLimit)
			}

			if requestNumber%10 == 0 {
				globalLimit *= 2
				readBufSize *= 2
				listener.SetConfig(globalLimit, connLimit)
			}
		}(conn)
	}
}
```
