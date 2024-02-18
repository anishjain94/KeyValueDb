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
	handleErr(err)
	defer l.Close()

	for {
		conn, err := l.Accept()
		handleErr(err)

		go processConn(conn)
	}

}

func processConn(conn net.Conn) {
	buf := make([]byte, 1024)

	for {
		read, err := conn.Read(buf)
		if err != nil {
			break
		}

		fmt.Println(string(buf[:read]))

		cms := parseCommand(string(buf[:read]))

		respond(conn, cms)

	}
}

func respond(conn net.Conn, cmds []string) {

	switch strings.ToLower(cmds[0]) {
	case "ping":
		conn.Write([]byte("+PONG\r\n"))

	case "echo":

		if len(cmds) < 2 {
			handleErr(fmt.Errorf("Invalid Args"))
		}
		fmt.Printf("$%d\r\n%s\r\n", len(cmds[1]), cmds[1])
		conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(cmds[1]), cmds[1])))

	}

}

func parseCommand(str string) []string {

	if len(str) == 0 {
		handleErr(fmt.Errorf("empty string"))
	}
	splitStr := strings.Split(str, "\r\n")

	numItems, err := strconv.Atoi(splitStr[0][1:])
	handleErr(err)

	neededItems := make([]string, numItems)

	index := 0
	for _, item := range splitStr {

		if len(item) > 0 && (item[0] != '$' && item[0] != '*') {
			neededItems[index] = item
			index++
		}
	}
	return neededItems
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
