package main

import (
	"fmt"
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
