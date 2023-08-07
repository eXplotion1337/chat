package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

var upgrader = websocket.Upgrader{}

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

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	// Получение текущей директории (где запущен файл main.go)
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}

	// Обслуживание статических файлов (index.html) с относительным путем
	parentDir := filepath.Dir(dir)
	webPath := filepath.Join(parentDir, "web")
	fs := http.FileServer(http.Dir(webPath))
	http.Handle("/", fs)


	fmt.Println("WebSocket server started. Listening on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
