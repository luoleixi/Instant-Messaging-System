package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// 创建用户的api
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	//启动监听当前消息的进程

	go user.ListenMessage()
	return user
}

// 用户上线功能
func (user *User) Online() {
	//用户上线，将用户加入到onlinMap中
	user.server.mapLock.Lock() //
	user.server.OnLineMap[user.Name] = user
	user.server.mapLock.Unlock()
	//广播当前用户上线消息
	user.server.BroadCat(user, "已上线")
}

// 用户下线功能
func (user *User) Offline() {
	//用户上线，将用户加入到onlinMap中
	user.server.mapLock.Lock() //
	delete(user.server.OnLineMap, user.Name)
	user.server.mapLock.Unlock()
	//广播当前用户上线消息
	user.server.BroadCat(user, "已下线")
}

// 用户传输消息
func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

// 用户处理消息功能
func (user *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户有哪些
		user.server.mapLock.Lock()
		for _, users := range user.server.OnLineMap {

			onlineMsg := "[" + users.Addr + "]" + users.Name + "在线"
			user.SendMsg(onlineMsg)
		}

		user.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		//判断用户名是否存在
		_, ok := user.server.OnLineMap[newName]
		if ok {
			user.SendMsg("当前用户名已存在")
		}
	} else {
		user.server.BroadCat(user, msg)
	}
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
