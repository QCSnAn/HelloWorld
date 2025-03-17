//go:build linux

package linuxonly

import (
	"context"
	"fmt"
	"net"
	"sync"
	"syscall"
	"test/testwork"
	"time"
	"unsafe"
)

func GetsockoptTCPInfo(tcpConn *net.TCPConn) (*syscall.TCPInfo, error) {
	file, err := tcpConn.File()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fd := file.Fd()
	tcpInfo := syscall.TCPInfo{}
	size := unsafe.Sizeof(tcpInfo)
	_, _, errno := syscall.Syscall6(syscall.SYS_GETSOCKOPT, fd, syscall.SOL_TCP, syscall.TCP_INFO,
		uintptr(unsafe.Pointer(&tcpInfo)), uintptr(unsafe.Pointer(&size)), 0)
	if errno != 0 {
		return nil, fmt.Errorf("syscall failed. errno=%d", errno)
	}

	return &tcpInfo, nil
}

func isConnected(connMap *sync.Map) {
	start := time.Now()

	tcpConnKey := make([]*net.TCPConn, 0)
	connMap.Range(func(key, value any) bool {
		tcpConn := key.(*net.TCPConn)
		tcpInfo, err := GetsockoptTCPInfo(tcpConn)
		fmt.Printf("TCP state is %d\n", tcpInfo.State)
		if err != nil {
			fmt.Println("Failed to get TCP info:", err)
			tcpConnKey = append(tcpConnKey, tcpConn)
			return true
		}
		switch tcpInfo.State {
		case testwork.TCP_LAST_ACK, testwork.TCP_CLOSE, testwork.TCP_FIN_WAIT1, testwork.TCP_FIN_WAIT2, testwork.TCP_TIME_WAIT:
			tcpConnKey = append(tcpConnKey, tcpConn)
			return true
		default:
			return true
		}
	})

	elapsed := time.Since(start)
	fmt.Println("The check took ", elapsed)

	for _, key := range tcpConnKey {
		value, ok := connMap.Load(key)
		if ok {
			value.(context.CancelFunc)()
			connMap.Delete(key)
		} else {
			fmt.Println("Key not found")
		}
	}
}
