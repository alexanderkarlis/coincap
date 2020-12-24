package table

import (
	"log"
	"sort"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// CoinList is an array of Coins.
// TODO: only need a type here??
type CoinList struct {
	Coins []Coin
}

// Coin struct
type Coin struct {
	Title string
	Value string
}

type CoinTable struct {
	wg    *sync.WaitGroup
	Table *widgets.Table
	Coins chan map[string]string
}

// NewTable func returns a new table instance
func NewTable(cs []string) *CoinTable {
	table := widgets.NewTable()
	table.Title = "💰 CoinCap 💰 (in USD)"

	styleMap := map[int]ui.Style{
		0: {Bg: ui.ColorMagenta},
	}

	initRows := func(c []string) [][]string {
		var vals [][]string
		var header, values []string

		for _, v := range c {
			header = append(header, v)
			values = append(values, "")
		}
		vals = append(vals, header)
		vals = append(vals, values)
		return vals
	}(cs)

	table.Rows = initRows
	table.RowStyles = styleMap
	width, height := ui.TerminalDimensions()
	table.SetRect(0, 0, width, height)

	t := &CoinTable{
		wg:    &sync.WaitGroup{},
		Table: table,
		Coins: make(chan map[string]string),
	}
	log.Println("starting table poll events...")
	t.wg.Add(1)
	go t.PollTableEvents()
	// ui.Render(table)
	return t
}

func (ct *CoinTable) RenderCoinTable() {
	ui.Render(ct.Table)
}

func (ct *CoinTable) PollTableEvents() {
	uiEvents := ui.PollEvents()
	defer ct.wg.Done()
	defer ui.Close()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				ct.Table.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(ct.Table)
			default:
			}
		case cs := <-ct.Coins:
			slice := mapTo2DSlice(cs)
			ct.Table.Rows = slice
			ui.Render(ct.Table)
		}
	}
}

func mapTo2DSlice(m map[string]string) [][]string {
	var vals [][]string
	var header, values []string

	m = sortMap(m)
	for k, v := range m {
		header = append(header, k)
		values = append(values, v)
	}
	vals = append(vals, header)
	vals = append(vals, values)
	return vals
}

func sortMap(m map[string]string) map[string]string {
	keys := make([]string, 0, len(m))
	newMap := make(map[string]string)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		newMap[k] = m[k]
	}

	return newMap
}
