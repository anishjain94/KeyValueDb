package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	handleErr(err)
	defer l.Close()

	for {

		fmt.Println("inside for")

		conn, err := l.Accept()
		handleErr(err)
		// defer conn.Close()

		b := make([]byte, 1024)
		n, err := conn.Read(b)
		handleErr(err)
		fmt.Println("number of bytes read", n)

		str := string(b)
		fmt.Println(str)
		count := strings.Count("PING", str)

		var output string
		for i := 0; i < count; i++ {
			output += "+PONG\r\n"
		}

		n, err = conn.Write([]byte(output))
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
