package main

import (
	"bufio"
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
		go processConn(conn)
	}
}

func processConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for true {
		cmds, err := readCmd(reader)
		if err != nil {
			return
		}

		cmds[0] = strings.ToLower(cmds[0])
		switch {
		case cmds[0] == "ping":
			cmdPing(conn)
		case cmds[0] == "echo":
			cmdEcho(conn, cmds)
		}
	}
}

func readCmd(reader *bufio.Reader) ([]string, error) {
	first, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	cmdCount, _ := strconv.Atoi(string(first)[1:])

	var cmds = make([]string, 0)
	for i := 0; i < cmdCount; i++ {
		_, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		buffer, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, string(buffer))
	}
	return cmds, err
}

func cmdPing(conn net.Conn) {
	fmt.Println("cmd ping.")
	response := "PONG"
	sendResponse(conn, response)
}

func cmdEcho(conn net.Conn, cmds []string) {
	fmt.Println("cmd echo.")
	response := strings.Join(cmds[1:], " ")
	sendResponse(conn, response)
}

func sendResponse(conn net.Conn, msg string) {
	length := len(msg)
	responseMsg := []byte("$" + strconv.Itoa(length) + "\r\n" + msg + "\r\n")
	conn.Write(responseMsg)
}
