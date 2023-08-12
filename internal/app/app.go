package app

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

var upgrader = websocket.Upgrader{}

func Run(config *Config) error {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	fs := http.FileServer(http.Dir(WEB_PATH))
	http.Handle("/", fs)

	fmt.Printf("WebSocket server started. Listening on %s\n", config.HttpConfing.Port)
	err := http.ListenAndServe(config.HttpConfing.Port, nil)
	if err != nil {
		return err
	}

	return nil
}

func handleRegister(w http.ResponseWriter, r *http.Request) {

}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		// Чтение данных из веб-сокета
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(clients, conn)
			break
		}
		// Преобразование полученных данных в строку и отправка в канал broadcast
		broadcast <- string(msg)
	}
}

func handleMessages() {
	for {
		// Получение сообщения из канала broadcast
		msg := <-broadcast
		// Отправка сообщения всем подключенным клиентам
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Error sending message:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
