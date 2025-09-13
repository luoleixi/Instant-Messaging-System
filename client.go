package  main

import (
	"fmt"
	"net"
	"flag"
	"io"
	"os"
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

//处理server回应的消息，直接显示道标准输出即可
func (client *Client) DealResponse() {
	io.Copy(os.Stdout,client.conn)
	//一旦client.conn 有数据，就直接copy道stdout 标准输出上，永久阻塞监听
	//for {
	//	buf := make()
	//	clinet.conn.Read(buf)
	//	fmt.Println(buf)
	//}
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

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_,err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:",err)
		return
	}
}

func (clinet *Client) PrivateChat() {
	var remoteName string
	var chatMsg string
	clinet.SelectUsers()
	fmt.Println(">>>>>请输入聊天对象[用户名],exit 退出<<<<<")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>请输入消息内容，exit 退出<<<<<")
		fmt.Scanln(&chatMsg)

		for chatMsg!= "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_,err := clinet.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:",err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>>请输入消息内容，exit 退出<<<<<")
			fmt.Scanln(&chatMsg)
		}

		clinet.SelectUsers()
		remoteName = ""
		fmt.Println(">>>>>请输入聊天对象[用户名],exit 退出<<<<<")
		fmt.Scanln(&remoteName)
	}

}

func (clinet *Client) PublicChat (){
	var chatMsg string
	//提示用户输入消息
	fmt.Println(">>>>>请输入聊天内容，exit退出<<<<<")
	fmt.Scanln(&chatMsg)

	//发给服务器
	for chatMsg != "exit" {
		//消息不为空则发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_,err := clinet.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:",err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>请输入聊天内容，exit退出<<<<<")
		fmt.Scanln(&chatMsg)
 	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>请输入用户名<<<<<")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_,err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:",err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1 :
			// 公聊模式
			client.PublicChat()
			break
		case 2 :
			// 私聊模式
			client.PrivateChat()
			break
		case 3 :
			// 更新用户名
			client.UpdateName()
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

	//单独开启一个go routine 去处理 server 的回执消息
	go client.DealResponse()

	fmt.Println (">>>>>连接服务器成功<<<<<")
	client.Run()
}