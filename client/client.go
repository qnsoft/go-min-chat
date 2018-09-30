package main

import (
	"net"
	"fmt"
	"os"
	"sync"
	"bufio"
	"flag"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
}

var pre string

func main() {
	var port string
	var host string
	flag.StringVar(&port, "p", "8080", "port")
	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.Parse()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	checkError(err)
	defer conn.Close()
	pre = fmt.Sprintf("%s:%s> ", host, port)
	fmt.Print(pre)
	var wg sync.WaitGroup
	ch := make(chan []byte)
	wg.Add(3)
	go readFromStdio(ch)
	go readFromConn(conn)
	go sendMsg(conn, ch)
	wg.Wait()
}

func readFromStdio(ch chan []byte) {
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		ch <- data
	}
}

func readFromConn(conn net.Conn) {
	for {
		buf := make([]byte, 50)
		n, err := conn.Read(buf)
		checkError(err)
		fmt.Print(string(buf[:n]))
		fmt.Println()
		fmt.Print(pre)
	}
}

func sendMsg(conn net.Conn, ch chan []byte) {
	for {
		content, _ := <-ch
		_, err := conn.Write(content)
		checkError(err)
		fmt.Print(pre)
	}
}
