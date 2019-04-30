package main

import (
	"fmt"
	"net"
	"os"
)

func Funcerr(s string,err error)  {
	if err!=nil{
		fmt.Println(s,err)
		os.Exit(-1)
	}
}
//创建一个用于维护在线用户的map
var OnLineUser map[string]User
//创建一个用于存放用户基本信息的结构体
type User struct {
	name string
	//添加一个用于接收广播消息的通道
	C chan string
}
//创建一个用于广播的通道
var Broadcast chan string
func main() {
	//创建监听套接字设置ip和端口
	listener,err:=net.Listen("tcp","127.0.0.1:8001")
	Funcerr("NetListen err:",err)
	defer listener.Close()
	//初始化广播通道
	Broadcast=make(chan string)
	//初始化map
	OnLineUser=make(map[string]User)
	//创建go程广播客户端发送的消息
	go WriteforClient()
	//循环监听客户端连接
	for{
		conn,err:=listener.Accept()
		Funcerr("Accept err:",err)
		//创建go程与客户端通信
		go WithClient(conn)


	}


}
//创建一个用于与客户端通信的函数
func WithClient(conn net.Conn){
	defer conn.Close()
	//初始化结构体
	var People User
	People.C=make(chan string)
	//获取客户端IP地址
	addr:=conn.RemoteAddr().String()
	People.name = "$ "+ addr
	buf:=make([]byte,4096)
	//添加在线用户信息
	OnLineUser[addr]=People
	//广播用户上线消息
	Broadcast<-addr+People.name+"上线了\n"
	//接收广播消息
	go func(){
		for msg:=range People.C{
			conn.Write([]byte(msg))
		}

	}()
	//循环读取客户端发送的消息
	for {
		n,_:=conn.Read(buf)
		Broadcast<-addr+People.name+" :"+string(buf[:n])+"\n"

	}

}
//创建广播客户端发送的消息
func WriteforClient(){

	//循环读取信息并发送给所有在线用户
	for {
		str:=<-Broadcast
		for _,word:=range OnLineUser{
			word.C <-str
		}
	}
}