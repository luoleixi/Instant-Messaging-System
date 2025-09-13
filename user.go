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
			user.SendMsg(onlineMsg + "\n")
		}

		user.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		//判断用户名是否存在
		_, ok := user.server.OnLineMap[newName]
		if ok {
			user.SendMsg("当前用户名已存在")
		} else {
			user.server.mapLock.Lock()
			delete(user.server.OnLineMap, user.Name)
			user.server.OnLineMap[newName] = user
			user.server.mapLock.Unlock()

			user.Name = newName
			user.SendMsg("您已更新用户名" + user.Name + "\n")

		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式 ： to|张三|消息内容

		//1 获取对方的用户名
		remoteName := strings.Split(msg,"|")[1]
		if remoteName == "" {
			user.SendMsg("消息格式不正确，请使用 \"to|name|msg\" 格式 \n ")
			return
		}
		//2 根据用户名 得到对方 User 对象
		remoteUser, ok := user.server.OnLineMap[remoteName]
		if !ok {
			user.SendMsg("该用户名不存在 \n")
			return
		}
		//3 获取消息内容 ， 通过对方的User对象将消息内容发送过去
		content := strings.Split(msg,"|")[2]
		if content == "" {
			user.SendMsg("无消息内容，请重发 \n")
			return
		}
		remoteUser.SendMsg(user.Name + ":" + content + "\n")
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
