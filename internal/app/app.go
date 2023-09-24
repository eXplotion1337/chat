package app

import (
	"chat/internal/app/usecase"
	"chat/internal/registry"
	"chat/internal/server/middlewares"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

var upgrader = websocket.Upgrader{}

var container *registry.Container

type credentials struct {
	Username string `json:"regUsername"`
	Password string `json:"regPassword"`
}

type loginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Run(config *Config) error {
	cnt, err := NewContainer(config)
	mux := http.NewServeMux()
	connHandler := http.HandlerFunc(handleConnections)
	container = cnt
	if err != nil {
		panic(err)
	}
	authMiddleware := middlewares.NewAuthMiddleware(container.AuthService)
	mux.HandleFunc("/signup", handleRegister)
	mux.HandleFunc("/signin", handleLogin)
	mux.Handle("/ws", authMiddleware(connHandler))
	//http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	fs := http.FileServer(http.Dir(WEB_PATH))
	mux.Handle("/", fs)

	fmt.Printf("WebSocket server started. Listening on %s", config.HttpConfing.Port)
	err = http.ListenAndServe(config.HttpConfing.Port, mux)
	if err != nil {
		return err
	}

	return nil
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var sf loginCredentials
	err := json.NewDecoder(r.Body).Decode(&sf)
	if err != nil {
		panic(err)
	}

	q := &usecase.UserLoginInput{
		Username: sf.Username,
		Password: sf.Password,
	}
	session, err := container.LoginUsecase.ProcessAuth(q)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:    "chat_uid",
		Value:   string(session.Token),
		Expires: session.ExpiresAt,
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("success"))
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var sf credentials
	err := json.NewDecoder(r.Body).Decode(&sf)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	input := usecase.UserRegisterInput{
		Username: sf.Username,
		Password: sf.Password,
	}
	err = container.RegisterUsecase.ProcessRegister(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
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
