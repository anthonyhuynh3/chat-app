package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
    "sync"
    "./utils"
)

func main() {
    // Create a listener on a specific address and port
    listener, err := net.Listen("tcp", "localhost:12345")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer listener.Close()

    clients := make(map[string]net.Conn)
    var mutex sync.Mutex

    fmt.Println("Server is listening for incoming connections...")

    for {
        // Accept incoming connections
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        // Handle client     
        go handleClient(conn,clients,&mutex)
    }
}

func handleClient(conn net.Conn, clients map[string]net.Conn, mutex *sync.Mutex) {
    defer conn.Close()

    // Read the username from the client
    reader := bufio.NewReader(conn)
    username, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading username:", err)
        return
    }
    username = strings.TrimSpace(username)

    // Store the client's connection in the map
    mutex.Lock()
    clients[username] = conn
    mutex.Unlock()

    fmt.Printf("User '%s' connected from %s\n", username, conn.RemoteAddr())

    // Handle messages from the client
    for {
        buffer := make([]byte, 1024)
        n, err := conn.Read(buffer)
        if err != nil {
            fmt.Println("Client connection closed.")
            mutex.Lock()
            delete(clients,username)
            mutex.Unlock()
            return
        }

        message := string(buffer[:n])
        status := utils.ProcessMessage(message,conn,clients,mutex)

        _, WriteErr := conn.Write([]byte(status))
        
        if WriteErr != nil {
            fmt.Println("Error sending message:", err)
            return
        }

    }
}
