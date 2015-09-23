package main

import (
	"fmt"
	"net"
	"log"
	"os"
	"bufio"
	"strconv"
)

const (
	defaultHost = "localhost"
	defaultPort = 9999
)

var (
	trace *log.Logger = log.New(os.Stdout, "TRACE: ", log.Ldate | log.Ltime | log.Lshortfile)
)

// To test your server implementation, you might find it helpful to implement a
// simple 'client runner' program. The program could be very simple, as long as
// it is able to connect with and send messages to your server and is able to
// read and print out the server's echoed response to standard output. Whether or
// not you add any code to this file will not affect your grade.
func main() {
	conn, err := net.Dial("tcp", defaultHost + ":" + strconv.Itoa(defaultPort))
	defer conn.Close()
	if err != nil {
		trace.Println("Couldn't dial ", err)
		os.Exit(1)
	}
	msg := "Hello, World"
	fmt.Fprintf(conn, msg + "\n")

	res, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print(res)
}
