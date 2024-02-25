package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestApp(t *testing.T) {

	str := "*2\r\n$4\r\necho\r\n$3\r\nhey\r\n"
	splitStr := strings.Split(str, "\r\n")

	numItems, err := strconv.Atoi(splitStr[0][1:])

	fmt.Println(numItems, err)
}

func TestParseCommand(t *testing.T) {

	// Set
	str := "*3\r\n$3\r\nset\r\n$5\r\nhello\r\n$6\r\nhorses\r\n"

	// Get
	str = "*2\r\n$3\r\nget\r\n$3\r\nhello\r\n"

	cmd, _ := parseCommand([]byte(str))

	fmt.Println(cmd)
}

func TestReadRdbFile(t *testing.T) {

	file, err := os.ReadFile("../dump.rdb")
	handleErr(err)

	// dst := make([]byte, hex.DecodedLen(len(file)))
	hx, _ := hex.DecodeString(string(file))
	fmt.Println(string(hx))
}
