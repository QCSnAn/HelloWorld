package testwork

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// 服务器地址
	serverAddress := "localhost:8080"

	// 连接到服务器
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("Hello World"))
	if err != nil {
		fmt.Println("failed to write to connection", err)
		return
	}
}
