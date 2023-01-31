package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	fmt.Println("Server is starting...")
	for true {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		fmt.Println("Accept a connection: " + conn.RemoteAddr().String())
		go readCmd(conn)
	}

}

func readCmd(conn net.Conn) {
	for true {
		var cmdByte = make([]byte, 2048)

		length, err := conn.Read(cmdByte)
		if err != nil && length == 0 {
			fmt.Println("Error reading streaming or len is 0: ", err.Error())
			break
		}
		cmd := strings.ToLower(string(cmdByte))
		lines := strings.Split(cmd, "\n")
		fmt.Println("receive len: " + strconv.Itoa(length))
		for i := 0; i < len(lines); i++ {
			msg := strings.TrimSpace(lines[i])
			fmt.Println("receive msg: " + msg)
			switch {
			case msg == "ping":
				readPingCmd(conn)
			}
		}

	}
}

func readPingCmd(conn net.Conn) {
	fmt.Println("Response PONG.")
	response := "PONG"
	length := len(response)
	responseMsg := []byte("$" + strconv.Itoa(length) + "\r\n" + response + "\r\n")
	conn.Write(responseMsg)
}
