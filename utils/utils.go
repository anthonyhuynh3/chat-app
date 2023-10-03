package utils

import (
    "strings"
    "net"
    "sync"
    "strconv"
    "fmt"
)

//Process client commands
func ProcessMessage(message string,conn net.Conn, clients map[string]net.Conn, mutex *sync.Mutex) string {
    messageTokens := strings.Split(message, " ")
    command := messageTokens[0]

    switch command {
    case "SEND":
        sendStatus := sendMessage(conn,clients,messageTokens[1], messageTokens[2])
        if(sendStatus != nil){
            return "ERROR"
        } else {
            return "\nMESSAGE SENT TO " + messageTokens[1]
        }
    case "USERS":
        return "USERS SUCCESFUL"
    case "HELP":
        return "HELP"
    case "LIST":
        return ListClients(clients)
    default:
        return "ERROR"
    }

    return "DEFAULT"
}

func ListClients(clients map[string]net.Conn) string {
    var allClients = "\nCurrent Logged in Users:\n"
    index := 1;

    for i := range clients {
        allClients += strconv.Itoa(index) + ". "  + i + "\n"
        index++
    }

    return allClients
}

func sendMessage(conn net.Conn, clients map[string]net.Conn, recipient string, message string) error {
    if message == "" || recipient == "" {
        return fmt.Errorf("Invalid message or recipient")
    }

    recipientConn, ok := clients[recipient]
    if !ok {
        return fmt.Errorf("Recipient '%s' not found", recipient)
    }

    _, err := recipientConn.Write([]byte("\n"+message))
    if err != nil {
        return err
    }

    return nil
}