package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var dest string
var verbosity int

func init() {
	flag.StringVar(&dest, "dest", "localhost", "host:port to connect")
	flag.IntVar(&verbosity, "verbosity", 0, "Verbosity level")
	flag.Parse()
}

func SogPrintln(a ...interface{}) (n int, err error) {
	if verbosity > 0 {
		return fmt.Println(a)
	}
	return 0, nil
}

func main() {
	SogPrintln(os.Args[0])

	_, err := net.ResolveTCPAddr("tcp", dest)
	if err != nil {
		fmt.Println("net.ResolveTCPAddr", err)
		return
	}

	SogPrintln(dest)

	var conn, err2 = net.Dial("tcp", dest)
	if err2 != nil {
		fmt.Println("net.Dial error:", err2)
	}

	SogPrintln(conn.LocalAddr(), " -> ", conn.RemoteAddr())

	loop(conn)
}

func loop(conn net.Conn) {
	errCh := make(chan error)
	readCh := make(chan int)

	errStd := make(chan error)
	readStd := make(chan int)

	var b [512]byte
	var bStd [512]byte
	go handleRead(conn, errCh, readCh, b[0:])
	go handleRead(os.Stdin, errStd, readStd, bStd[0:])

	for {
		select {
		case err := <-errCh:
			SogPrintln("Error Conn", err)
			conn.Close()
			time.Sleep(10)
			return
		case n := <-readCh:
			SogPrintln("Read", n)
			fmt.Print("Read ", string(b[0:n]))
			go os.Stdout.Write(b[0:n])
		case err := <-errStd:
			SogPrintln("Error Std", err)
			os.Stdin.Close()
			time.Sleep(10)
			return
		case n := <-readStd:
			SogPrintln("Read", n)
			fmt.Print("Read ", string(bStd[0:n]))
			go handleWrite(conn, errCh, bStd[0:n])
		}
	}
}

type SogConn interface {
	Read(b []byte) (n int, err error)

	Write(b []byte) (n int, err error)

	Close() error
}

func handleWrite(conn SogConn, errCh chan error, b []byte) {
	tot := len(b)
	n, err := conn.Write(b)
	if err != nil {
		errCh <- err
		if n != tot {
			SogPrintln("Short write", n, "of", tot, "bytes")
		}
	}
}

func handleRead(conn SogConn, errCh chan error, ch chan int, b []byte) {

	for {
		n, err := conn.Read(b)
		if err != nil {
			errCh <- err
			break
		}

		ch <- n
	}
}

func toBytes(s string) []byte {
	return []byte(s)
}
