package testwork

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"
)

const (
	TCP_ESTABLISHED = 1
	TCP_SYN_SENT    = 2
	TCP_SYN_RECV    = 3
	TCP_FIN_WAIT1   = 4
	TCP_FIN_WAIT2   = 5
	TCP_TIME_WAIT   = 6
	TCP_CLOSE       = 7
	TCP_CLOSE_WAIT  = 8
	TCP_LAST_ACK    = 9
	TCP_LISTEN      = 10
)

const checkInterval = 1

var connMap sync.Map

func checkConnected() {

	ticker := time.Tick(time.Minute)

	for {
		select {
		case <-ticker:
			fmt.Println("check func existing...")
		default:
		}
		isConnected()
		time.Sleep(checkInterval * time.Second)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		fmt.Println("Error casting to TCPConn")
		return
	}
	ctx, cancel := context.WithCancel(context.TODO())
	connMap.Store(tcpConn, cancel)
	defer func() {
		_, ok = connMap.Load(tcpConn)
		if ok {
			connMap.Delete(tcpConn)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task interruption")
			return
		default:
			buf := make([]byte, 1024)
			length, err := conn.Read(buf)
			if err != nil {
				fmt.Println("failed to read from connection", err)
				return
			}

			for i := 0; i < 10; i++ {
				select {
				case <-ctx.Done():
					fmt.Println("Task interruption")
					return
				default:
					fmt.Println("Do something...")
					time.Sleep(10 * time.Second)
				}
			}

			_, err = conn.Write(buf[:length])
			if err != nil {
				fmt.Println("failed to write to connection", err)
				return
			}
		}
	}
}

func startAccept(ln net.Listener) {
	go checkConnected()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}
		go handleConn(conn)
	}
}

func main() {
	TCPAddr := "localhost:8080"
	//UNIXAddr := "file/unixsock.sock"

	TCPln, err := net.Listen("tcp", TCPAddr)
	if err != nil {
		fmt.Println("Error Creating TCPln:", err)
		return
	}
	defer TCPln.Close()

	//UNIXln, err := net.Listen("unix", UNIXAddr)
	//if err != nil {
	//	fmt.Println("Error Creating UNIXln:", err)
	//	return
	//}
	//defer UNIXln.Close()

	go startAccept(TCPln)
	//go startAccept(UNIXln)

	for {
		numGoroutines := runtime.NumGoroutine()
		fmt.Printf("Number of goroutines: %d\n", numGoroutines)
		time.Sleep(1 * time.Second)
	}
}
