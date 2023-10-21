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

func echo(l net.Listener) {
	log.Println("TCP Echo Start: ", l.Addr())
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}

			go func(conn net.Conn) {
				n, err := io.Copy(conn, conn)
				log.Printf(
					"TCP Echo From: %s Echo: %d, Error: %v",
					conn.RemoteAddr(), n, err,
				)
				conn.Close()
			}(c)
		}
	}()
}

func blackhole(l net.Listener) {
	log.Println("TCP Blackhole Start: ", l.Addr())
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}

			go func(conn net.Conn) {
				go func() {
					var s int64
					buf := make([]byte, 32768)

					for {
						n, err := conn.Write(buf)
						if err != nil {
							break
						}
						s += int64(n)
					}
					log.Printf(
						"TCP Blackhole Send To: %s Size: %d, Error: %v",
						conn.RemoteAddr(), s, err,
					)
				}()

				n, err := io.Copy(io.Discard, conn)
				log.Printf(
					"TCP Blackhole Receive From: %s Blackhole: %d, Error: %v",
					conn.RemoteAddr(), n, err,
				)

				conn.Close()
			}(c)
		}
	}()

}

func main() {
	var port string
	var addr string
	var mode string

	flag.StringVar(&addr, "addr", "127.0.0.1", "Blackhole TCP Address")
	flag.StringVar(&port, "port", "9999", "Blackhole TCP Port")
	flag.StringVar(&mode, "mode", "blackhole", "blackhole for Blackhole Server, echo for Echo Server")
	flag.Parse()

	l, err := net.Listen("tcp", net.JoinHostPort(addr, port))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer l.Close()
	switch mode {
	case "blackhole":
		blackhole(l)
	case "echo":
		echo(l)
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("TCP Blackhole Shutdown.")
}
