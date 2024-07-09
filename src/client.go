package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIP string
	SeverPort int
	Name string
	conn net.Conn
}

func NewClient(ip string, port int) *Client {
	// Create Client Obj
	clientP := &Client{
		ServerIP: ip,
		SeverPort: port,
	}
	// Connect to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return nil
	}

	clientP.conn = conn

	// Return obj
	return clientP
}

// Init before main
var serverIP string
var serverPort int

func init(){
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "Set IP of Server (127.0.0.1 by default)")
	flag.IntVar(&serverPort, "port", 8888, "Set Port of Server (8888 by default)")
}

func main() {
	// Parsing Command Lines
	flag.Parse()

	// Create Client
	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>> Connect to Server unsuccessfully...")
	} else {
		fmt.Println(">>>>>>>> Connect to Server succesfully...")
	}
	
	// Launch Client Business
	select{}
}