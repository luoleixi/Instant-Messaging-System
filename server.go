package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnLineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

// 创建server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnLineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听Message广播消息 channel
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, cli := range s.OnLineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// 广播消息的方法
func (s *Server) BroadCat(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	//fmt.Println("链接创建成")

	user := NewUser(conn, s)

	//用户上线，将用户加入到onlinMap中
	user.Online()
	//接受客户端发送的消息

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			//提取用户消息
			msg := string(buf[:n-1])
			//消息进行广播
			user.DoMessage(msg)
		}
	}()

	select {}
}

// 启动服务器的接口
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	go s.ListenMessager()
	//accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		//do handler
		go s.Handler(conn)
	}
}
