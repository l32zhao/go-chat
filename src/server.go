package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
)

type Server struct {
	Ip		string
	Port	int

	// Online Users Map
	OnlineMap map[string]*User
	mapLock sync.RWMutex

	// Message Channel for BroadCasting
	MC chan string
}

// Server port
func NewServer(ip string, port int) *Server {
	serverP := &Server{
		Ip:		ip,
		Port:	port,
		OnlineMap: make(map[string]*User),
		MC: make(chan string),
	}

	return serverP
}


// BroadCasting
func (this *Server) ListenUserMsg() {
	for {
		msg := <-this.MC

		// Send to each online users
		this.mapLock.Lock()
		for _, oc := range this.OnlineMap {
			oc.C <-msg
		}
		this.mapLock.Unlock()

	}
}
func (this *Server) BroadCast(user *User, msg string) {
	infodMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.MC <-infodMsg + "\n"
}

// Handler
func (this *Server) Handler(conn net.Conn) {
	// Bussiness Logic... to be connected
	fmt.Println("Successfully Connected")

	user := NewUser(conn, this)

	// Make user Online
	user.Online()

	// Alive
	isAlive := make(chan bool)

	// Recv msg from clients
	go func() {
		buf := make([]byte, 4096)
		for{
			n, err := conn.Read(buf)
			// Offline if empty msg
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read failed:", err)
				return
			}

			// Extract message and Remove empty line
			msg := string(buf[:n-1])

			// Handle specific user's msg
			user.HandleMsg(msg)

			// It's alive if receiving any msg
			isAlive <- true
		}
	}()

	// Tracking user status
	const timeout = 60
	for {
		// Blocking current handler
		select {
		case <- isAlive:
			// Do nothing, the following statements would still be runned
			// Although only one case (isAlive) will return to select
		case <- time.After(time.Second * timeout):	// Timeout
			// Terminate this user
			user.SendMsg("You are out!")
			// Close
			close(user.C)
			conn.Close()
	
			// Exit: return or runtime.Goexit()
			runtime.Goexit()
		}
	}
}

// Launch server
func (this *Server) Start() {
	// Socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("Listening failed:", err)
		return
	}

	// Close
	defer listener.Close()

	// Launch goroutine for Message
	go this.ListenUserMsg()
	
	for {
		// Accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener accept failed:", conn)
			return
		}
		// Do handler
		go this.Handler(conn)
	}
	
}