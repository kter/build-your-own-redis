package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

func reply(conn net.Conn, replyString string) {
	data := []byte(replyString)
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println("cannot write", err)
	}
}

func chop(s string) string {
	s = strings.TrimRight(s, "\n")
	if strings.HasSuffix(s, "\r") {
		s = strings.TrimRight(s, "\r")
	}

	return s
}
func processUndefinedCommand(command string) (string, error) {
	return "", errors.New("Undefined Command: " + command)
}

func processPingCommand() (string, error) {
	return "+PONG\r\n", nil
}

func processCommand(command string) (string, error) {
	switch chop(command) {
	case "PING":
		return processPingCommand()
	default:
		return processUndefinedCommand(command)
	}
}

func getRequestString(conn net.Conn, c chan string) {
	data := make([]byte, 1024)
	count, _ := conn.Read(data)
	c <- string(data[:count])
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		fmt.Println("Connection Accepted!")
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		for {
			c := make(chan string)
			go getRequestString(conn, c)
			receivedCommand := <-c
			fmt.Println(receivedCommand)

			replyString, err := processCommand(receivedCommand)
			if err != nil {
				replyString = "Error accepting connection: " + err.Error()
			}

			reply(conn, replyString)
		}
	}
}
