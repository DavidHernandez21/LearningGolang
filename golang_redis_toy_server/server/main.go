package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type command int

const (
	Get  command = iota + 1 // 1
	Set                     // 2
	Incr                    // 3
	Del                     // 4
)

type commandMessage struct {
	commandName     command
	key             string
	value           string
	responseChannel chan string
}

func preProcessInput(input string) []string {

	trimmedInput := strings.TrimSpace(input)
	if trimmedInput == "" {
		return []string{}
	}
	partsRaw := strings.Split(trimmedInput, " ")
	parts := make([]string, len(partsRaw))
	newLength := 0
	for i := range partsRaw {
		if partsRaw[i] == "" {
			continue
		}
		parts[newLength] = partsRaw[i]
		newLength++
	}

	returnValue := make([]string, newLength)
	_ = copy(returnValue, parts[:newLength])

	return returnValue
}

func handleConnection(commandChannel chan<- commandMessage, conn net.Conn, getChannel chan string, setChannel chan string, incrChannel chan string, deleteChannel chan string) {
	defer conn.Close()

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("error reading:", err)
			// if errors.Is(err, os.ErrDeadlineExceeded) {
			// 	fmt.Println("Deadline exceeded")
			// 	_, err = conn.Write([]byte("Dealine exceeded" + "\n"))
			// 	if err != nil {
			// 		fmt.Println("error writing:", err)
			// 		return
			// 	}
			// 	return
			// }
		}

		var response string
		parts := preProcessInput(netData)

		if len(parts) == 0 {
			_, err = conn.Write([]byte("No input was provided" + "\n"))
			if err != nil {
				fmt.Println("error writing:", err)
				return
			}
			continue
		}
		command := parts[0]
		// getChannel := make(chan string)
		// setChannel := make(chan string)
		// incrChannel := make(chan string)
		// deleteChannel := make(chan string)

		switch command {
		case "STOP", "QUIT":
			return
		case "GET":
			response = handleGet(commandChannel, getChannel, parts)
		case "SET":
			response = handleSet(commandChannel, setChannel, parts)
		case "INCR":
			response = handleIncrement(commandChannel, incrChannel, parts)
		case "DEL":
			response = handleDelete(commandChannel, deleteChannel, parts)
		default:
			response = "ERR unknown command"
		}

		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			fmt.Println("error writing:", err)
			return
		}
	}
}

func handleGet(commandChannel chan<- commandMessage, getChannel chan string, commandSlice []string) string {
	if len(commandSlice) < 2 {

		return "ERR wrong number of arguments for 'get' command"

	}
	commandMessage := commandMessage{
		commandName:     Get,
		key:             commandSlice[1],
		responseChannel: getChannel}

	commandChannel <- commandMessage
	return <-commandMessage.responseChannel

}

func handleSet(commandChannel chan<- commandMessage, setChannel chan string, commandSlice []string) string {
	if len(commandSlice) < 3 {

		return "ERR wrong number of arguments for 'set' command"

	}

	commandMessage := commandMessage{
		commandName:     Set,
		key:             commandSlice[1],
		value:           commandSlice[2],
		responseChannel: setChannel}

	commandChannel <- commandMessage
	return <-commandMessage.responseChannel

}

func handleIncrement(commandChannel chan<- commandMessage, incrChannel chan string, commandSlice []string) string {
	if len(commandSlice) < 2 {

		return "ERR wrong number of arguments for 'incr' command"

	}
	commandMessage := commandMessage{
		commandName:     Incr,
		key:             commandSlice[1],
		responseChannel: incrChannel}

	commandChannel <- commandMessage
	return <-commandMessage.responseChannel
}

func handleDelete(commandChannel chan<- commandMessage, delChannel chan string, commandSlice []string) string {
	if len(commandSlice) < 2 {

		return "ERR wrong number of arguments for 'del' command"

	}
	commandMessage := commandMessage{
		commandName:     Del,
		key:             commandSlice[1],
		responseChannel: delChannel}

	commandChannel <- commandMessage
	return <-commandMessage.responseChannel
}

func handleDB(commandChannel chan commandMessage) {

	db := make(map[string]string, 50)

	for command := range commandChannel {
		switch command.commandName {
		case Get:
			command.responseChannel <- db[command.key]
		case Set:
			db[command.key] = command.value
			command.responseChannel <- "OK"
		case Incr:
			handleDBIncrement(command, db)
		case Del:
			handleDBDelete(command, db)
		}
	}

}

func handleDBIncrement(command commandMessage, db map[string]string) {
	value, ok := db[command.key]
	var response string

	if !ok {
		response = "1"
		db[command.key] = response
		command.responseChannel <- response
		return
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		// response =
		command.responseChannel <- "ERR value is not an integer or out of range"
		return
	}
	response = strconv.Itoa(intValue + 1)
	db[command.key] = response
	command.responseChannel <- response

}

func handleDBDelete(command commandMessage, db map[string]string) {
	_, ok := db[command.key]

	if !ok {
		command.responseChannel <- "0"
		return

	}
	delete(db, command.key)
	command.responseChannel <- command.key
}

func signalHandler(server net.Listener, commandChannel chan commandMessage) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	fmt.Printf("Recieved signal %q\n", <-c)
	fmt.Println("Shutting down...")
	close(commandChannel)
	if err := server.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)

}

func main() {

	const PORT = ":8080"

	server, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	println("Server started on port", PORT)

	commandChannel := make(chan commandMessage)

	defer func() {
		fmt.Printf("Closing server on port %s", PORT)
		close(commandChannel)
		if err := server.Close(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	go signalHandler(server, commandChannel)

	go handleDB(commandChannel)

	getChannel := make(chan string)
	setChannel := make(chan string)
	incrChannel := make(chan string)
	deleteChannel := make(chan string)

	// ticker := time.NewTicker(3 * time.Second)
	for {
		// select {
		// case <-ticker.C:
		// 	fmt.Printf("30 seconds passed, closing server on port %s", PORT)
		// 	return

		conn, err := server.Accept()
		// conn.SetReadDeadline(time.Now().Add(15 * time.Second))
		// conn.SetWriteDeadline(time.Now().Add(20 * time.Second))
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(commandChannel, conn, getChannel, setChannel, incrChannel, deleteChannel)
	}

}
