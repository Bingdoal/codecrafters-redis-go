package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var storageMap = map[string]string{}
var timerMap = map[string]*time.Timer{}

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
		switch cmds[0] {
		case "ping":
			cmdPing(conn)
		case "echo":
			cmdEcho(conn, cmds)
		case "set":
			cmdSet(conn, cmds)
		case "get":
			cmdGet(conn, cmds)
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

func cmdSet(conn net.Conn, cmds []string) {
	fmt.Println("cmd set")
	key := cmds[1]
	storageMap[key] = cmds[2]

	if len(cmds) >= 4 &&
		strings.ToLower(cmds[3]) == "px" {
		expiredTime := cmds[4]
		expired, err := strconv.Atoi(expiredTime)
		if err != nil {
			fmt.Errorf("expired: %s error: %s", expiredTime, err.Error())
			return
		}

		_, containTimer := timerMap[key]
		if containTimer && expired >= 0 {
			timerMap[key].Reset(time.Duration(expired) * time.Second)
		} else if containTimer && expired < 0 {
			timerMap[key].Stop()
			delete(timerMap, key)
		} else if !containTimer && expired >= 0 {
			timer := time.AfterFunc(time.Duration(expired)*time.Millisecond, func() {
				delete(storageMap, key)
				delete(timerMap, key)
			})
			timerMap[key] = timer
		}
	}
	sendResponse(conn, "OK")
}

func cmdGet(conn net.Conn, cmds []string) {
	fmt.Println("cmd get")
	value, contain := storageMap[cmds[1]]
	if !contain {
		sendNullResponse(conn)
	} else {
		sendResponse(conn, value)
	}
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
	completeMsg := "$" + strconv.Itoa(length) + "\r\n" + msg + "\r\n"
	responseMsg := []byte(completeMsg)
	conn.Write(responseMsg)
}

func sendNullResponse(conn net.Conn) {
	completeMsg := "$-1\r\n\r\n"
	responseMsg := []byte(completeMsg)
	conn.Write(responseMsg)
}
