package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"
)

//global variables
var (
	//tcp
	ports          []string   //ports keep the clients port connection
	connections    []net.Conn //keeps every connection
	newConnections []net.Conn //auxiliary slice

	//rpc
	serverPorts []string //keeps the list of ports used for the chat rooms
	chatTopics  []string //keeps chats topics
	chatUsers   []int64  //keeps the number of users connected per chat room
)

func server(port, thematic string) {
	//listen to client connection on a defined port
	s, err := net.Listen("tcp", port)
	sub := strings.Split(port, ":")
	port = sub[1]
	ports = append(ports, port)

	if err != nil {
		fmt.Println(err)
		return
	}
	//keeps on listening
	for {
		c, err := s.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}
		connections = append(connections, c)
		go handleClient(c, port)
	}
}

func handleClient(c net.Conn, port string) {
	var msg string
	for {
		err := gob.NewDecoder(c).Decode(&msg)

		if err != nil {
			fmt.Println(err)
			return
		}

		//when the message contains that, it means that it's their first connection
		//we add the username to the slice
		if strings.Contains(msg, ":") {
			receive := strings.Split(msg, ":")
			ports = append(ports, receive[1])
			switch port {
			case "9997":
				chatUsers[0]++
			case "9999":
				chatUsers[2]++
			case "9998":
				chatUsers[1]++
			}
		}

		//if it's not a disconnection we send the message to all clients
		//excluding the sender
		if msg != "disconnect" {
			for i := 0; i < len(connections); i++ {
				if c != connections[i] && ports[i] == port && !strings.Contains(msg, ":") {
					err := gob.NewEncoder(connections[i]).Encode(msg)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		} else { //if the message is disconnect we close the connection
			//and we erase it from the slice
			for i := 0; i < len(connections); i++ {
				if c == connections[i] {
					c.Close()
					fmt.Println("Client has disconnected")
				} else {
					newConnections = append(newConnections, connections[i])
				}
			}
			connections = newConnections
			newConnections = nil
		}
	}
}

type Server struct{}

func (this *Server) GetPorts(name string, reply *[]string) error {
	*reply = ports
	return nil
}
func (this *Server) GetChatRoomsTopics(name string, reply *[]string) error {
	*reply = chatTopics
	return nil
}
func (this *Server) GetChatRoomsUsers(name string, reply *[]int64) error {
	*reply = chatUsers
	return nil
}

//rpc Server registering
func rpcServer() {
	rpc.Register(new(Server))
	ln, err := net.Listen("tcp", ":9996")
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
	var thematic string

	scanner := bufio.NewScanner(os.Stdin)

	//tcp connections to clients
	//chat rooms
	fmt.Println("Enter the first chat room thematic")
	scanner.Scan()
	thematic = scanner.Text()
	chatTopics = append(chatTopics, thematic)
	chatUsers = append(chatUsers, 0)
	go server(":9997", thematic)

	fmt.Println("Enter the second chat room thematic")
	scanner.Scan()
	thematic = scanner.Text()
	chatTopics = append(chatTopics, thematic)
	chatUsers = append(chatUsers, 0)
	go server(":9998", thematic)

	fmt.Println("Enter the third chat room thematic")
	scanner.Scan()
	thematic = scanner.Text()
	chatTopics = append(chatTopics, thematic)
	chatUsers = append(chatUsers, 0)
	go server(":9999", thematic)

	//rpc server
	go rpcServer()

	//holding server running
	fmt.Println("\nStarting server...")
	var input string
	fmt.Scanln(&input)
}
