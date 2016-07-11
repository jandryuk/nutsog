package main

import (
	"fmt"
	"flag"
	"os"
	"net"
)

var port int
var dest string

func init() {
	flag.StringVar(&dest, "dest", "localhost", "host:port to connect")
	flag.IntVar(&port, "port", 8080, "port number to connect to")
	flag.Parse()
}

func main() {
	var name string
	name = os.Args[0]
	fmt.Println(name)

	net.LookupIP(dest)

	var conn, err = net.Dial("tcp", dest)
	if (err != nil) {
		fmt.Println("net.Dial error:", err)
	}

	fmt.Println(conn)

	fmt.Println("LocalAddr", conn.LocalAddr())
	fmt.Println("RemoteAddr", conn.RemoteAddr())

	var n int
	n, err = conn.Write(toBytes("Happy\n"))
	if (err != nil) {
		fmt.Println("conn.Write error:", err)
	}

	fmt.Println("Wrote", n, "bytes")
}

func toBytes(s string) []byte {
	return []byte(s)
}
