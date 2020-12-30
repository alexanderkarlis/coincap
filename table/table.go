package table

import (
	"fmt"
	"sort"
	"strconv"
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

// CoinTab is the tab for each coin passed through command line
type CoinTab struct {
	wg              *sync.WaitGroup
	mux             *sync.Mutex
	TabPane         *widgets.TabPane
	Coins           chan map[string]string
	CoinValuesArray map[string][]float64
	tabsToCoins     map[int]string
}

// NewTabs Funtions
func NewTabs(cs []string) *CoinTab {
	tabpane := widgets.NewTabPane(cs...)
	mm := make(map[string][]float64)
	ttc := make(map[int]string)
	for _, coin := range cs {
		mm[coin] = []float64{}
	}
	for i := 0; i < len(cs); i++ {
		ttc[i] = cs[i]
	}

	tabpane.Border = true
	tabpane.SetRect(0, 1, 50, 4)

	t := &CoinTab{
		wg:              &sync.WaitGroup{},
		mux:             &sync.Mutex{},
		TabPane:         tabpane,
		Coins:           make(chan map[string]string),
		CoinValuesArray: mm,
		tabsToCoins:     ttc,
	}
	t.wg.Add(1)

	// start the service for events to the table
	go t.PollTabEvents()
	return t
}

func (ct *CoinTab) getActiveCoinTab() string {
	return ct.tabsToCoins[ct.TabPane.ActiveTabIndex]
}

func (ct *CoinTab) renderTabPane() {
	activeCoin := ct.getActiveCoinTab()
	cts := ct.CoinValuesArray[activeCoin][1:]

	p1 := widgets.NewPlot()
	p1.Marker = widgets.MarkerBraille
	p1.SetRect(0, 10, 100, 50)
	p1.AxesColor = ui.ColorWhite
	p1.LineColors[0] = ui.ColorYellow
	p1.DrawDirection = widgets.DrawLeft
	p1.PlotType = widgets.LineChart
	p1.HorizontalScale = 4
	p1.BorderStyle.Fg = ui.ColorGreen
	p1.Title = activeCoin

	p2 := widgets.NewParagraph()
	p2.Title = "current price ($)"
	p2.SetRect(0, 5, 20, 10)
	p2.BorderStyle.Fg = ui.ColorCyan

	p3 := widgets.NewParagraph()
	p3.Title = "% Change"
	p3.SetRect(25, 5, 45, 10)

	if len(cts) > 2 {
		p2.Text = fmt.Sprintf("%f", cts[len(cts)-1])
		pChange := 100 * (cts[len(cts)-1] - cts[len(cts)-2]) / cts[len(cts)-2]

		if pChange > 0 {
			p3.TextStyle.Fg = ui.ColorGreen
		} else if pChange <= 0 {
			p3.TextStyle.Fg = ui.ColorRed
		}
		p3.Text = fmt.Sprintf("%f %s", pChange, "%")

		p1.Data = make([][]float64, 1)
		p1.Data = [][]float64{cts}
		p1.MaxVal = cts[len(cts)-1] + cts[len(cts)-1]*0.1

		// p1.Max.Y = int(cts[len(cts)-1] + cts[len(cts)-1]*0.1)
		// p1.Min.Y = int(cts[len(cts)-1] - cts[len(cts)-1]*0.1)

		ui.Render(ct.TabPane, p3, p2, p1)
	}
}

// PollTabEvents continually polls for events to the widget to trigger a redraw
// I.e. - UI input or data input.
func (ct *CoinTab) PollTabEvents() {
	fmt.Println()
	uiEvents := ui.PollEvents()
	defer ct.wg.Done()
	defer ui.Close()
	ui.Render(ct.TabPane)
	// ct.renderTabPane()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				// os.Exit(1)
				return
			case "h":
				ct.TabPane.FocusLeft()
				ui.Clear()
				ui.Render(ct.TabPane)
				ct.renderTabPane()
			case "l":
				ct.TabPane.FocusRight()
				ui.Clear()
				ui.Render(ct.TabPane)
				ct.renderTabPane()
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				ct.TabPane.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ct.renderTabPane()
			default:
			}
		case coinData := <-ct.Coins:
			for key, val := range coinData {
				fl, _ := strconv.ParseFloat(val, 8)
				ct.mux.Lock()
				ct.CoinValuesArray[key] = append(ct.CoinValuesArray[key], fl)
				ct.mux.Unlock()
			}
			ct.renderTabPane()
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
