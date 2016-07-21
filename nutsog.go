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

func StderrPrintln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stderr, a)
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
		return
	}

	SogPrintln(conn.LocalAddr(), " -> ", conn.RemoteAddr())

	loop(conn)
}

func loop(conn net.Conn) {
	errNet := make(chan error)
	netCh := make(chan []byte)

	errStd := make(chan error)
	stdCh := make(chan []byte)

	go handleReadBuf(os.Stdin, conn, errStd, stdCh)
	go handleReadBuf(conn, os.Stdin, errNet, netCh)

	var netTot, stdTot int

	for {
		select {
		case err := <-errNet:
			SogPrintln("Error Conn", err)
			conn.Close()
			time.Sleep(10)
			return
		case buf := <-netCh:
			netTot += len(buf)
			StderrPrintln("NetRead", len(buf), netTot)
			go handleWrite(os.Stdout, errStd, buf)
		case err := <-errStd:
			SogPrintln("Error Std", err)
			os.Stdin.Close()
			time.Sleep(10)
			return
		case buf := <-stdCh:
			stdTot += len(buf)
			StderrPrintln("StdRead", len(buf), stdTot)
			go handleWrite(conn, errNet, buf)
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

func handleReadBuf(conn SogConn, wconn SogConn, errCh chan error, ch chan []byte) {

	var b [512]byte

	for {
		n, err := conn.Read(b[0:])
		if err != nil {
			errCh <- err
			break
		}

		buf := make([]byte, n)
		copy(buf, b[:n])

		handleWrite(wconn, errCh, buf)
	}
}

func toBytes(s string) []byte {
	return []byte(s)
}
