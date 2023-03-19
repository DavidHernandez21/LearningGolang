package main

import (
	"bufio"
	"net"
	"os"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			println("Close failed:", err.Error())
			os.Exit(1)
		}
	}()

	for {

		var reader *bufio.Reader = bufio.NewReader(os.Stdin)

		println("Enter text to send to server: ")

		text, err := reader.ReadString('\n')

		if err != nil {
			println("ReadString failed:", err.Error())
			break
		}

		if strings.Contains(text, "exit") {
			println("Exiting TCP client!")
			break
		}

		_, err = conn.Write([]byte(text))
		if err != nil {
			println("Write data failed:", err.Error())
			break
		}

		// buffer to get data
		received := make([]byte, 1024)
		_, err = conn.Read(received)
		if err != nil {
			println("Read data failed:", err.Error())
			break
		}

		println("Received message:", string(received))

	}
}
