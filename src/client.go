package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIP string
	SeverPort int
	Name string
	conn net.Conn

	modeFlag int	// Mode
}

func NewClient(ip string, port int) *Client {
	// Create Client Obj
	clientP := &Client{
		ServerIP: ip,
		SeverPort: port,
		modeFlag: 999,	// Prevent from exiting immediately
	}
	// Connect to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("Failed to connect to server: ", err)
		return nil
	}

	clientP.conn = conn

	// Return obj
	return clientP
}

// Query Online User
func (client *Client) queryUsers() {
	sendMsg := "?\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("Write err:", err)
	}
	return
}

// Mode
func (client *Client) privateChat() {
	// Prompt
	var userName string
	var chatMsg string

	client.queryUsers()
	fmt.Println(">>>>>>Enter a user name, type exit() to exit>>>>>>>")
	fmt.Scanln(&userName)

	for userName != "exit()" {
		fmt.Println(">>>>>>Start Chatting... type exit() to exit>>>>>>>")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit()" {
			// Send msg
			if len(chatMsg) != 0 {
				// send through 'to|username|msg'
				sendMsg := "-to|" + userName + "|" + chatMsg + "\n\n"	// double \n for new line
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("Sending failed:", err)
					break
				}
			}
			
			chatMsg = ""
			fmt.Println(">>>>>>Start Chatting... type exit() to exit>>>>>>>")
			fmt.Scanln(&chatMsg)
		}
		client.queryUsers()
		fmt.Println(">>>>>>Enter a user name, type exit() to exit>>>>>>>")
		fmt.Scanln(&userName)
	}
}

func (client *Client) publicChat() {
	// Prompt
	var chatMsg string

	fmt.Println(">>>>>>Please Chat with others, type exit() to exit>>>>>>>")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit()" {
		// Send msg
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("Sending failed:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>>Start Chatting... type exit() to exit>>>>>>>")
		fmt.Scanln(&chatMsg)
	}
}


func (client *Client) rename() bool {
	fmt.Println(">>>>>>>Please input User Name:")
	fmt.Scanln(&client.Name)

	msg := "-r " + client.Name + "\n"
	_, err := client.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("Failed to write:", err)
		return false
	}

	return true
}

// Handle Response from Server
func (client *Client) handleResponse() {
	io.Copy(os.Stdout, client.conn)	// Blocking forever to copy byte from stream stdout

	// Equals to this
	// for {
	// 	buf := make([]byte, 1024)
	// 	client.conn.Read(buf)
	// 	fmt.Println(string(buf))
	// }
}

// Run client based on diff modes
func (client *Client) mode() bool {
	var flag int

	fmt.Println("Enter 1 for Public Mode")
	fmt.Println("Enter 2 for Private Mode")
	fmt.Println("Enter 3 for Updating User Name")
	fmt.Println("Enter 0 to Exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.modeFlag = flag
		return true
	}
	fmt.Println(">>>>>>>>Please Enter Legal Number>>>>>>>>>")
	return false
}

func (client *Client) Run() {
	for client.modeFlag != 0 {
		// Keep asking if illegal input
		for client.mode() != true {}

		// Switch to diff business logic
		switch client.modeFlag {
		case 1:
			// fmt.Println("Public Mode...")
			client.publicChat()
			break
		case 2:
			// fmt.Println("Private Mode...")
			client.privateChat()
			break
		case 3:
			// fmt.Println("Rename Mode...")
			client.rename()
			break
		}
	}
	return
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
	
	go client.handleResponse()

	// Launch Client Business
	client.Run()
}