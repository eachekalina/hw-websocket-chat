package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/", echo)
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		return
	}
}

var clients = make(map[*websocket.Conn]struct{})

func echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()

	clients[connection] = struct{}{}  // Сохраняем соединение, используя его как ключ
	defer delete(clients, connection) // Удаляем соединение

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
		}

		// Теперь мы рассылаем сообщения всем клиентам
		go writeMessage(message)

		go messageHandler(message)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

func writeMessage(message []byte) {
	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}

func messageHandler(message []byte) {
	fmt.Println(string(message))
}
