package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var port string
	var addr string

	flag.StringVar(&addr, "addr", "127.0.0.1", "Blackhole TCP Address")
	flag.StringVar(&port, "port", "9999", "Blackhole TCP Port")
	flag.Parse()

	l, err := net.Listen("tcp", net.JoinHostPort(addr, port))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer l.Close()
	log.Println("TCP Blackhole Start: ", l.Addr())
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}

			go func(conn net.Conn) {
				n, err := io.Copy(io.Discard, conn)
				log.Printf(
					"TCP Blackhole Receive From: %s Blackhole: %d, Error: %v",
					conn.RemoteAddr(), n, err,
				)
				conn.Close()
			}(c)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("TCP Blackhole Shutdown.")
}
