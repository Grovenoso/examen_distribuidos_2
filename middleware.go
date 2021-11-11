package main

import (
	"fmt"
	"net"
	"net/rpc"
)

////MIDDLEWARE - SERVER
func ConnectServer(conn *rpc.Client, err error) []string {
	var result []string
	err = conn.Call("Server.GetChatRoomsInfo", result, &result)
	if err != nil {
		fmt.Println(err)
	} else {
		return result
	}
	return nil
}

func getServerIP() {

}

////MIDDLEWARE - CLIENT
type Middleware struct{}

func (this *Middleware) GetServerInfo(name string, reply *[]string) error {
	conn, err := rpc.Dial("tcp", ":9996")
	if err != nil {
		fmt.Println(err)
		return err
	}
	var result []string
	err = conn.Call("Server.GetChatRoomsInfo", name, &result)
	if err != nil {
		fmt.Println(err)
	} else {
		*reply = result
	}
	return nil
}

func rpcMiddleware() {
	rpc.Register(new(Middleware))
	//middleware - client connection on port 9995
	ln, err := net.Listen("tcp", ":9995")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func main() {
	go rpcMiddleware()

	//holding middleware running
	fmt.Println("\nStarting middleware...")
	var input string
	fmt.Scanln(&input)
}
