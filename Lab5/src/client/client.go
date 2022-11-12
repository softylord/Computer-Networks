package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	//"os/signal"

	"github.com/gorilla/websocket"
)

/*Вход: Массы и координаты
центра масс пяти тел
в трехмерном
пространстве

Выход: Координата центра
масс системы*/

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	fmt.Println("Массы и координаты 5 тела вводятся в одну строчку (сначала масса, потом координаты),\nпараметры каждого тела обособляются скобками")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		vectors := scanner.Text()
		switch vectors {
		case "quit":
			os.Exit(0)
		default:
			err := c.WriteMessage(websocket.TextMessage, []byte(vectors))
			if err != nil {
				log.Println("write:", err)
				return
			}
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("Центр масс системы: %s", message)

		}
	}
}
