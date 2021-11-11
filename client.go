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

func getServerInfo() []string {
	conn, err := rpc.Dial("tcp", ":9995")
	if err != nil {
		fmt.Println(err)
	}
	var result []string
	var name string
	err = conn.Call("Middleware.GetServerInfo", name, &result)
	if err != nil {
		fmt.Println(err)
	} else {
		return result
	}
	return nil
}

func main() {
	var opc int
	//var port string
	//var status = make(chan int)
	//var msg = make(chan string)

	fmt.Println("Enter your username")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userName := scanner.Text()

	fmt.Println("Welcome, " + userName + "!")
	fmt.Println("Choose a chatroom from the ones below...")

	options := getServerInfo()

	for i := 0; i < len(options); i++ {
		fmt.Print(i + 1)
		fmt.Println(". " + options[i])
	}
	fmt.Println("Or enter 0 to exit")

	//keep listening to user input
	for {
		fmt.Scanln(&opc)
		switch opc {
		case 1:
			fmt.Println("You chose chatroom 1: ", options[0])
		case 2:
			fmt.Println("You chose chatroom 2: ", options[1])
		case 3:
			fmt.Println("You chose chatroom 3: ", options[2])
		case 0:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Wrong input")
		}
	}
}
