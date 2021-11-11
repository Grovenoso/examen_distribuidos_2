package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"net/rpc"
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
				err = gob.NewEncoder(c).Encode(port)
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

func getServerTopics() []string {
	conn, err := rpc.Dial("tcp", ":9995")
	if err != nil {
		fmt.Println(err)
	}
	var result []string
	var name string
	err = conn.Call("Middleware.GetServerTopics", name, &result)
	if err != nil {
		fmt.Println(err)
	} else {
		return result
	}
	return nil
}

func getServerUsers() []int {
	conn, err := rpc.Dial("tcp", ":9995")
	if err != nil {
		fmt.Println(err)
	}
	var result []int
	var name string
	err = conn.Call("Middleware.GetServerUsers", name, &result)
	if err != nil {
		fmt.Println(err)
	} else {
		return result
	}
	return nil
}

func connectClientToServer(option int, userName string) {
	var port string
	var optChat int
	var status = make(chan int)
	var msg = make(chan string)

	switch option {
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

	//Chatroom menu
	fmt.Println("\nMenu" +
		"\n 1. Send message" +
		"\n 2. Stop client")

	//keep listening to user input
	for {
		fmt.Scanln(&optChat)
		switch optChat {
		//send message
		case 1:
			fmt.Println("Enter your message")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			message := scanner.Text()
			message = userName + " - " + message
			status <- 1
			msg <- message
		//disconnect
		case 2:
			status <- 2
			msg <- "disconnect"
			return
		//any other input won't be accepted
		default:
			fmt.Println("\nWrong input")
		}
	}
}

func main() {
	var opt int

	fmt.Println("Enter your username")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userName := scanner.Text()

	fmt.Println("Welcome, " + userName + "!")
	fmt.Println("Choose a chatroom from the ones below..." + "\n")

	options := getServerTopics()
	users := getServerUsers()

	for i := 0; i < len(options); i++ {
		fmt.Print(i + 1)
		fmt.Print(". " + options[i] + " - ")
		fmt.Println(users[i])
	}
	fmt.Println("Or enter 0 to exit")

	//keep listening to user input
	for {
		fmt.Scanln(&opt)
		switch opt {
		case 1:
			fmt.Println("You chose chatroom 1: ", options[0])
			connectClientToServer(opt, userName)

		case 2:
			fmt.Println("You chose chatroom 2: ", options[1])
			connectClientToServer(opt, userName)

		case 3:
			fmt.Println("You chose chatroom 3: ", options[2])
			connectClientToServer(opt, userName)

		case 0:
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Wrong input")
		}
	}
}
