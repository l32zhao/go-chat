package main

import (
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

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>>>>> Connect to Server unsuccessfully...")
	} else {
		fmt.Println(">>>>>>>> Connect to Server succesfully...")
	}
	
	// Launch Client Business
	select{}
}