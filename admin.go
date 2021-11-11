package main

import (
	"fmt"
	"net/rpc"
)

func getServerPorts() []string {
	conn, err := rpc.Dial("tcp", ":9995")
	if err != nil {
		fmt.Println(err)
	}
	var result []string
	var name string
	err = conn.Call("Middleware.GetServerPorts", name, &result)
	if err != nil {
		fmt.Println(err)
	} else {
		return result
	}
	return nil
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

func updateServerState() {
	options := getServerTopics()
	users := getServerUsers()
	ports := getServerPorts()

	for i := 0; i < len(options); i++ {
		fmt.Print(ports[i] + ": " + options[i] + " - ")
		fmt.Println(users[i])
	}
}

func main() {
	var opt int

	fmt.Println("Welcome, admin!..." + "\n")

	updateServerState()

	//keep listening to admin input
	for {
		fmt.Println("\nMenu" +
			"\n 1. Update chatrooms states" +
			"\n 2. Stop admin")
		fmt.Scanln(&opt)
		switch opt {
		case 1:
			fmt.Println("Updating...")
			updateServerState()

		case 2:
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Wrong input")
		}
	}
}
