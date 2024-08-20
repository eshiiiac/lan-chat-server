package main

import (
	//"fmt"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10240,
	WriteBufferSize: 10240,
}

type Client struct {
	Conn *websocket.Conn
	IP   string
}

// array for clients
var clients []Client

func main() {

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		//extractting the clients' ip
		clientIP := r.RemoteAddr
		fmt.Println("New Client connected: ", clientIP)

		client := Client{
			Conn : conn,
			IP : clientIP,
		}

		clients = append(clients, client)

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			for _, client := range clients {
				messageWithIP := fmt.Sprintf("From %s: %s", clientIP, string(msg))
				if err = client.Conn.WriteMessage(msgType, []byte(messageWithIP)); err != nil {
					return
				}
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.ListenAndServe(":8080", nil)
}
