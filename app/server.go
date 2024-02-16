package main

import (
	"fmt"
	"log"
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
		defer conn.Close()

		b := make([]byte, 1024)
		n, err := conn.Read(b)

		n, err = conn.Write([]byte("+PONG\r\n"))
		handleErr(err)

		log.Printf("received %d bytes", n)
	}

}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
