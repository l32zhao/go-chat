package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C chan string	// channel fd
	conn net.Conn	// socket fd

	server *Server	// Inited in server and passed to user
}

// Create User
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()	// Get remote IP address

	userP := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,

		server: server,
	}

	// Create Goroutine for each user channel
	go userP.ListenMsg()

	return userP
}

// Listen current user receiving channel
func (this *User) ListenMsg() {
	for {
		msg := <-this.C
		
		this.SendMsg(msg)
	}
}

func (this *User) Online() {
	// Add user to OnlineMap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// BroadCasting to other users
	this.server.BroadCast(this, " now is online!")
}

func (this *User) Offline() {
	// Remove user from OnlineMap
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// BroadCasting to other users
	this.server.BroadCast(this, " now is terminated.")
}

// Send
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))	// Write to related user's client (not to server)
}

// Do indicated tasks
func (this *User) HandleMsg(msg string) {
	switch {
	case msg == "?":	// query current OnlineMap
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			infoMsg := "[" + user.Addr + "]" + user.Name + ":" + " is online...\n"
			this.SendMsg(infoMsg)
		}
		this.server.mapLock.Unlock()
	case len(msg) > 3 && msg[:3] == "-r ":
		// newName := strings.Split(msg, '|')[1]
		newName := msg[3:]

		_, exist := this.server.OnlineMap[newName]
		if exist {
			this.SendMsg("Rename Failed: Current user name is existed!\n")
		} else {
			// Update OnlineMap on Server
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)	// Without delete, two same users would exist
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			
			// Update Local Name
			this.Name = newName
			this.SendMsg("Your user name is renamed as: " + this.Name + "\n")
		}
	case len(msg) > 4 && msg[:4] == "-to|":
		// Get User Name
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("Empty Format: please use \"-to|name|msg\".\n")
			return
		}

		// Get User server
		remoteUser, exist := this.server.OnlineMap[remoteName]
		if !exist {
			this.SendMsg("User not exist!\n")
			return
		}

		// Get msg content
		context := strings.Split(msg, "|")[2]
		if context == "" {
			this.SendMsg("Empty msg, resend!\n")
			return
		}
		
		// Send privately
		remoteUser.SendMsg(this.Name + " said: " + context + "\n")

	default:
		// Basic Broadcasting
		this.server.BroadCast(this, msg)
	}
}