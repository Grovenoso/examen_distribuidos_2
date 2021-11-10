package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

//global variables
var (
	ports          []string   //ports keep the clients port connection
	connections    []net.Conn //keeps every connection
	newConnections []net.Conn //auxiliary slice
)

func server(port, thematic string) {
	//listen to client connection on a defined port
	s, err := net.Listen("tcp", port)
	sub := strings.Split(port, ":")
	port = sub[1]

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

		//when the message contains that, means that it's their first connection
		//we add the username to the slice
		if strings.Contains(msg, ":") {
			receive := strings.Split(msg, ":")
			ports = append(ports, receive[1])
		}

		//if it's not a disconnection we send the message to all clients
		//excluding the sender
		if msg != "disconnect" || !strings.Contains(msg, ":") {
			for i := 0; i < len(connections); i++ {
				fmt.Println(ports[i], port)
				if c != connections[i] && ports[i] == port {
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

func main() {
	var thematic string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter the first chat room thematic")
	scanner.Scan()
	thematic = scanner.Text()
	go server(":9999", thematic)

	fmt.Println("Enter the second chat room thematic")
	scanner.Scan()
	thematic = scanner.Text()
	go server(":9998", thematic)

	fmt.Println("Enter the third chat room thematic")
	scanner.Scan()
	thematic = scanner.Text()
	go server(":9997", thematic)

	fmt.Println("\nStarting server...")

	//goroutine server
	fmt.Scanln()
}
