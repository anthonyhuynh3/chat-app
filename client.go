package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	fmt.Print("Enter your username: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := scanner.Text()

	conn, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	// Send the username to the server
	_, err = fmt.Fprintf(conn, "%s\n", username)
	if err != nil {
		fmt.Println("Error sending username:", err)
		return
	}

	var wg sync.WaitGroup

	// Create a channel to receive messages from the server
	messageChan := make(chan string)

	// Goroutine to read and print messages from the server
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			buffer := make([]byte, 1024)
			n, ReadErr := conn.Read(buffer)
			if ReadErr != nil {
				fmt.Println("Server connection closed.")
				close(messageChan)
				return
			}
			message := string(buffer[:n])
			messageChan <- message
		}
	}()

	// Goroutine to read user input and send messages to the server
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			fmt.Print("Enter a message (or 'exit' to quit): ")
			scanner.Scan() // Read the user's input
			message := scanner.Text()
			message = strings.TrimSpace(message)

			if message == "exit" {
				fmt.Println("Exiting the program.")
				return
			}

			_, err := conn.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		}
	}()

	// Goroutine to print received messages
	wg.Add(1)
	go func() {
		defer wg.Done()
		for message := range messageChan {
			fmt.Println(message)
			fmt.Print("Enter a message (or 'exit' to quit): ")
		}
	}()

	wg.Wait()
}
