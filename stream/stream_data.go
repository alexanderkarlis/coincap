package stream

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/alexanderkarlis/coincap/rabbit"
	"github.com/gorilla/websocket"
)

var addr = flag.String("web", "ws.coincap.io", "coinbase websocket service address")
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Main is the entrypoint into the stream code -> RabbitMQ
func Main() {
	coins := os.Args[1:]
	Stream(coins)
}

// Stream function to get data from Coinbase Stream
func Stream(coinList []string) {
	coinString := []byte(strings.Join(coinList, ", "))
	rabbit.SendCoinNames(coinString)
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var queryString string

	queryString = strings.Join(coinList, ",")
	queryString = "assets=" + queryString

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/prices", RawQuery: queryString}
	// log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			go rabbit.SendPrices(message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
