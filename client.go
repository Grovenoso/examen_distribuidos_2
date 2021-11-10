package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

//dialing the server
func runClient(userName, port string, status chan int, msg chan string) {
	//connection on the chosen port
	c, err := net.Dial("tcp", port)

	if err != nil {
		fmt.Println(err)
		return
	}

	var message, received string

	//welcome message
	message = port

	//always listening
	go func() {
		defer c.Close()
		for {
			err := gob.NewDecoder(c).Decode(&received)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(received)
		}
	}()

	//keep the connection
	for {
		select {
		case _status := <-status:
			//first connection
			if _status == 0 {
				err = gob.NewEncoder(c).Encode(message)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			//send message
			if _status == 1 {
				message = <-msg
				err := gob.NewEncoder(c).Encode(message)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			//terminate connection
			if _status == 2 {
				message = <-msg
				err := gob.NewEncoder(c).Encode(message)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("\nGoodbye!")
				return
			}
		}
	}
}

func main() {
	var opc int
	var port string
	var status = make(chan int)
	var msg = make(chan string)

	//until the client enters their name
	//the client won't connect to the server
	fmt.Println("Enter your username")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userName := scanner.Text()

	//choose chat room
	fmt.Println("Choose chat room" +
		"\n1. Chat room 1" +
		"\n2. Chat room 2" +
		"\n3. Chat room 3")

	fmt.Scanln(&opc)
	switch opc {
	case 1:
		port = ":9997"
	case 2:
		port = ":9998"
	case 3:
		port = ":9999"
	}

	//goroutine connecting to server
	go runClient(userName, port, status, msg)
	status <- 0

	//main menu
	fmt.Println("\nMenu" +
		"\n 1. Send message" +
		"\n 2. Stop client")

	//keep listening to user input
	for {
		fmt.Scanln(&opc)

		switch opc {
		//send message
		case 1:
			fmt.Println("Enter your message")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			message := scanner.Text()
			message = userName + ": " + message
			status <- 1
			msg <- message
		//disconnect
		case 2:
			status <- 3
			msg <- "disconnect"
			return
		//any other input won't be accepted
		default:
			fmt.Println("\nWrong option")
		}
	}
}
