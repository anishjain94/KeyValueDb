package main

import (
	"fmt"
	"net"
	"os"
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
		_, err := conn.Read(buf)
		if err != nil {
			break
		}
		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Error responding to connection: ", err.Error())
			os.Exit(1)
		}

	}
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
