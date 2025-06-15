package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
	server *Server
}

// 创建用户的api
func NewUser(conn net.Conn,server *Server) *User {
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

//用户上线功能
func (user *User) Online(){
	//用户上线，将用户加入到onlinMap中
	user.server.mapLock.Lock() //
	user.server.OnLineMap[user.Name] = user
	user.server.mapLock.Unlock()
	//广播当前用户上线消息
	user.server.BroadCat(user, "已上线")
}

//用户下线功能
func(user *User) Offline(){
	//用户上线，将用户加入到onlinMap中
	user.server.mapLock.Lock() //
	delete(user.server.OnLineMap,user.Name)
	user.server.mapLock.Unlock()
	//广播当前用户上线消息
	user.server.BroadCat(user, "已下线")
}

//用户处理消息功能
func(user *User) DoMessage(msg string) {
	user.server.BroadCat(user,msg)
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
