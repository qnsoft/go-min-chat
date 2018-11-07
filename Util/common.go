package Util

import (
	"fmt"
	"os"
	"net"
	"io"
	"github.com/beego/bee/logger/colors"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
}

func ConnReadCheckError(err error, conn net.Conn) int {
	if err != nil {
		if (err == io.EOF) {
			fmt.Printf("client %s is close!\n", conn.RemoteAddr().String())
			return -1
		} else {
			fmt.Println("Error: %s", err.Error())
			return 0
		}
	}
	return 1
}

func EchoLine(content string, level int) {
	var s string
	switch level {
	case 1:
		s = colors.Green(content)
	case 2:
		s = colors.Red(content)
	}
	fmt.Println(s)
}