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

	conn, err := l.Accept()
	handleErr(err)
	b := make([]byte, 1024)

	for {
		n, err := conn.Read(b)
		handleErr(err)

		log.Printf("received %d bytes", n)
		log.Printf("received the following data: %s", string(b[:n]))

		// go func() {

		output := "+PONG\r\n"

		_, err = conn.Write([]byte(output))
		handleErr(err)
		// }()
	}

}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
