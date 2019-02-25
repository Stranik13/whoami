package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var port string

func init() {
	flag.StringVar(&port, "port", ":80", "give me a port number")
}

func main() {
	flag.Parse()

	fmt.Println("Starting up on port " + port)
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatal("error opening port: %v", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go serveTCP(conn)
	}
}

func serveTCP(conn net.Conn) {
	defer conn.Close()

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		} else if temp == "WHO" {
			result := whoAmIInfo()
			conn.Write([]byte(result))
		} else {
			result := fmt.Sprintf("Received: %s", netData)
			conn.Write([]byte(result))
		}
	}
}

func whoAmIInfo () string {
	var out bytes.Buffer

	hostname, _ := os.Hostname()
	out.WriteString(fmt.Sprintf("Hostname: %s\n", hostname))

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			out.WriteString(fmt.Sprintf("IP: %s\n", ip))
		}
	}

	return out.String()
}