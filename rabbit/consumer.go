package rabbit

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/alexanderkarlis/coincap/view"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/streadway/amqp"
)

var coinPricesMapChan chan map[string]string

var coinPricesMap map[string]string
var coinNames []string

var table *widgets.Table

func init() {
	coinPricesMap = make(map[string]string)
	coinPricesMapChan = make(chan map[string]string)
}

func cosumeCoinNames(channel *amqp.Channel) {
	q, err := channel.QueueDeclare(
		"coin-names", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	names, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	for d := range names {
		str := string(d.Body)
		strSplit := strings.Split(str, ", ")
		for i := range strSplit {
			coinNames = append(coinNames, strSplit[i])
		}
		coinNames = unique(coinNames) // Get unique values only
		cMap := make(map[string]string)
		for _, c := range coinNames {
			cMap[c] = "0"
		}
		coinPricesMap = cMap
		return
	}
}

// Consume , function to Consume data from producer
func Consume() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	table := widgets.NewTable()
	table.Title = "CoinCap coins 💰"
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	width, height := ui.TerminalDimensions()
	table.SetRect(0, 0, width, height)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	rabbitChannel, err := conn.Channel()
	defer rabbitChannel.Close()

	cosumeCoinNames(rabbitChannel)
	FailOnError(err, "Failed to open a channel")

	q, err := rabbitChannel.QueueDeclare(
		"coin-prices", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := rabbitChannel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			coinPricesMapChan <- convertToMap(d.Body)
		}
	}()

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				table.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(table)
			}
		case <-coinPricesMapChan:
			mapValue, ok := <-coinPricesMapChan
			if ok == false {
				fmt.Println(mapValue, ok, "loop is broken coinPricesMapChan")
				return
			}
			for colkey := range mapValue {
				coinPricesMap[colkey] = mapValue[colkey]
			}
			newRows := view.MapToRows(coinPricesMap)
			table.Rows = newRows
			ui.Render(table)
		default:

		}
	}
}

func convertToMap(byt []byte) map[string]string {
	var dat map[string]string
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}

	return dat
}

func unique(elements []string) []string {
	elementMap := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		elementMap[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range elementMap {
		result = append(result, key)
	}
	return result
}
