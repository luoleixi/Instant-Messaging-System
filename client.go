package  main

import (
	"fmt"
	"net"
	"flag"
)

type Client struct {
	ServerIp string  	//服务器 ip
	ServerPort int  	//服务器端口
	Name string 		//用户名称
	conn net.Conn 		//连接句柄
	flag int 			//当前client模式
}

func NewClient(serverIp string ,  serverPort int) *Client {
	//创建客户端对象
	clinet := &Client{
		ServerIp: serverIp,
		ServerPort :serverPort,
		flag : -1 ,
	}
	//链接server
	conn,err := net.Dial("tcp",fmt.Sprintf("%s:%d",serverIp,serverPort))
	if err != nil {
		fmt.Println("net.Dial err:",err)
		return nil
	}
	clinet.conn = conn
	//返回对象
	return clinet
}

func (clinet *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		clinet.flag = flag
		return true
	} else {
		fmt.Println(">>>>>请输入合法范围内的数字<<<<<")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1 :
			// 公聊模式
			fmt.Println("公聊模式选择...")
			break
		case 2 :
			// 私聊模式
			fmt.Println("私聊模式选择...")
			break
		case 3 :
			// 更新用户名
			fmt.Println("更新用户名选择...")
			break
		}
	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init () {
	flag.StringVar(&serverIp,"ip","127.0.0.1","设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort,"port",8888,"设置服务器端口(默认是8888)")
}

func main() {
	flag.Parse()

	client := NewClient(serverIp,serverPort)
	if client == nil {
		fmt.Println (">>>>>连接服务器失败<<<<<")
		return
	}
	fmt.Println (">>>>>连接服务器成功<<<<<")
	client.Run()
}