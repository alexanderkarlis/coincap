package streamer

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gorilla/websocket"

	"github.com/alexanderkarlis/coincap/table"
)

var (
	addr     = flag.String("web", "ws.coincap.io", "coinbase websocket service address")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type SocketMessage struct {
	wg          sync.WaitGroup
	mux         sync.RWMutex
	Message     string
	Coinlist    []string
	CoinMap     map[string]string
	Table       table.CoinTable
	quitChan    chan struct{}
	chanMessage chan map[string]string
}

func init() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
}

func NewStreamer(coinList []string) *SocketMessage {
	m := make(map[string]string)
	for i := 0; i < len(coinList); i++ {
		m[coinList[i]] = ""
	}
	ct := table.NewTable(coinList)
	s := SocketMessage{
		wg:          sync.WaitGroup{},
		mux:         sync.RWMutex{},
		Message:     "",
		Coinlist:    coinList,
		Table:       *ct,
		CoinMap:     m,
		quitChan:    make(chan struct{}),
		chanMessage: make(chan map[string]string),
	}
	s.wg.Add(1)
	go s.worker()

	return &s
}

func (s *SocketMessage) worker() {
	defer s.wg.Done()
	for {
		select {
		case <-s.quitChan:
			log.Println("Stop message received. Exiting worker.")
			return
		case coinmap := <-s.chanMessage:
			for k, v := range coinmap {
				if k == "" {
					continue
				}
				s.mux.Lock()
				s.CoinMap[k] = v
				s.mux.Unlock()
			}
			// add RW mutex here for the map string string
			// do I need a RLock here?
			s.mux.RLock()
			s.Table.Coins <- s.CoinMap
			s.mux.RUnlock()
		}
	}
}

func (s *SocketMessage) QueueCoinData(m map[string]string) bool {
	select {
	case s.chanMessage <- m:
		return true
	default:
		return false
	}
}

func (s *SocketMessage) GetCoinData() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	queryString := "assets=" + strings.Join(s.Coinlist, ",")

	// add url encoding of ws request
	u := url.URL{Scheme: "wss", Host: *addr, Path: "/prices", RawQuery: queryString}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	go func(sig chan os.Signal, doneCh chan struct{}) {
		for {
			select {
			case <-sig:
				doneCh <- struct{}{}
			}
		}
	}(interrupt, done)

	defer close(done)
	for {
		select {
		case <-done:
			os.Exit(1)
		default:
			m := make(map[string]string)
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read message error:", err)
				return
			}

			err = json.Unmarshal(message, &m)
			if err != nil {
				log.Fatal("could not unmarshal message")
			}
			queued := s.QueueCoinData(m)
			if !queued {
				os.Exit(1)
			}
		}
	}
}
