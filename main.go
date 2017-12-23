package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func showLocalAddrs() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		fmt.Println(addr.String())
	}
}

// Listen - receive function
func Listen(port int) error {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	defer lis.Close()
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("Accept error:", err)
		}
		log.Println("accept:", conn.RemoteAddr())

		go func(c net.Conn) {
			io.Copy(os.Stdout, c)
			log.Println("closed:", conn.RemoteAddr())
			defer c.Close()
		}(conn)
	}
}

// Dial - send function
func Dial(host string, port int) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	defer conn.Close()
	go io.Copy(os.Stdout, conn)
	_, err = io.Copy(conn, os.Stdin)
	return err
}

func main() {
	port := flag.Int("p", 0, "local port number")
	flag.Usage = func() {
		fmt.Println(strings.Replace(
			`options:
connect to somewhere:	$name hostname port
listen:			$name -p port
	-p		listen port number`,
			"$name", filepath.Base(os.Args[0]), -1))
	}
	flag.Parse()
	if *port > 0 {
		log.Fatal(Listen(*port))
	}

	if flag.NArg() != 2 {
		flag.Usage()
		return
	}
	dialPort := 0
	fmt.Sscanf(flag.Arg(1), "%d", &dialPort)
	if err := Dial(flag.Arg(0), dialPort); err != nil {
		log.Println(err)
	}
}
