package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	flag.StringVar(&Dir, "dir", "", "dir value")
	flag.StringVar(&DbFileName, "dbfilename", "", "dbname value")
	flag.StringVar(&port, "port", "6379", "alternate port")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	handleErr(err)
	defer l.Close()

	for {
		conn, err := l.Accept()
		handleErr(err)

		go processConn(conn)
	}

}

func processConn(conn net.Conn) {
	go HandleExpiry(expiryChannel, redis)
	buf := make([]byte, 2024)

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
		conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(commands.args[0]), commands.args[0])))

	case "set":
		redis[commands.args[0]] = RedisValues{
			Value: commands.args[1],
		}
		conn.Write([]byte("+OK\r\n"))

		if len(commands.args) == 4 {
			expiry, err := strconv.Atoi(commands.args[3])
			handleErr(err)

			redis[commands.args[0]] = RedisValues{
				Value: commands.args[1],
			}

			expiryChannel <- ExpiryStruct{Key: commands.args[0], Expiry: time.Duration(expiry) * time.Millisecond}

		}

	case "get":

		if val, ok := redis[commands.args[0]]; ok {

			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val.Value), val.Value)))
		} else {
			conn.Write([]byte("$-1\r\n"))
		}
	case "config":

		if commands.args[0] == "get" {
			if commands.args[1] == "dir" {
				// *2 $3 dir $16 /tmp/redis-files

				// *2\r\n$3\r\ndir\r\n$16\r\n/tmp/redis-files\r\n

				temp := fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(commands.args[1]), commands.args[1], len(Dir), Dir)
				conn.Write([]byte(temp))
			}

			if commands.args[1] == "dbfilename" {
				temp := fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(commands.args[1]), commands.args[1], len(DbFileName), DbFileName)
				conn.Write([]byte(temp))
			}
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
			args = append(args, strings.ToLower(strings.Trim(val, " ")))
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

func HandleExpiry(expiryChannel chan ExpiryStruct, redis map[string]RedisValues) {

	for val := range expiryChannel {
		tiker := time.NewTicker(val.Expiry)

		<-tiker.C
		fmt.Println("deleting ", val.Key)
		delete(redis, val.Key)
		tiker.Stop()

	}

}
