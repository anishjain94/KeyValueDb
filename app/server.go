package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
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
		cmds, err := parseCommand(buf[:read])
		handleErr(err)

		respond(conn, cmds)

	}
}

func respond(conn net.Conn, commands RedisCommand) {

	switch strings.ToLower(commands.command) {
	case "ping":
		conn.Write([]byte("+PONG\r\n"))

	case "echo":
		fmt.Printf("$%d\r\n%s\r\n", len(commands.args[0]), commands.args[0])
		conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(commands.args[0]), commands.args[0])))

	case "set":
		redis[commands.args[0]] = RedisValues{
			Value: commands.args[1],
		}
		conn.Write([]byte("+OK\r\n"))

		if len(commands.args) == 4 {
			expiry, err := strconv.Atoi(commands.args[3])
			handleErr(err)

			expiredTime := time.Now().Add(time.Duration(expiry) * time.Millisecond)
			redis[commands.args[0]] = RedisValues{
				Value:     commands.args[1],
				ExpiresAt: &expiredTime,
			}

		}

	case "get":

		if val, ok := redis[commands.args[0]]; ok {

			if val.ExpiresAt != nil && time.Now().After(*val.ExpiresAt) {
				conn.Write([]byte("$-1\r\n"))
			}

			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val.Value), val.Value)))
		} else {
			conn.Write([]byte("$-1\r\n"))
		}

	}

}

func parseCommand(buffer []byte) (RedisCommand, error) {
	input := strings.Fields(string(buffer))
	result := RedisCommand{}

	argLen, err := strconv.Atoi(input[0][1:])
	handleErr(err)

	result.command = strings.ToUpper(input[2])
	if argLen <= 1 {
		return result, nil
	}
	args := make([]string, 0, argLen-1)
	for _, val := range input[4:] {
		if !strings.HasPrefix(val, "$") && len(val) > 0 {
			args = append(args, strings.Trim(val, " "))
		}
	}
	result.args = args
	return result, nil

}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
