package main

import (
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8085", Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	i := 0
	for _ = range ticker.C {
		err := c.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(i)))
		i++
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
