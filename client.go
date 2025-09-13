package  main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp string  	//服务器 ip
	ServerPort int  	//服务器端口
	Name string 		//用户名称
	conn net.Conn 		//连接句柄
}

func NewClient(serverIp string ,  serverPort int) *Client {
	//创建客户端对象
	clinet := &Client{
		ServerIp: serverIp,
		ServerPort :serverPort,
	}
	//链接server
	conn,err := net.Dial("tcp",fmt.Sprintf("%s:%d",serverIp,serverPort))
	if err != nll {
		fmt.Println("net.Dial err:",err)
		return nil
	}

	clinet.conn = conn

	//返回对象
	return clinet

}

func main() {
	client := NewClient("127.0.0.1",8888)
	if client == nil {
		fmt.Println (">>>>>连接服务器失败<<<<<")
		return
	}
	fmt.Println (">>>>>连接服务器成功<<<<<")
	select {

	}
}