package view

import (
	"fmt"
	"sort"

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

var msg chan []byte
var coinNames []string
var startRows chan [][]string
var previousValue chan [][]string

func init() {
	startRows = make(chan [][]string)
	previousValue = make(chan [][]string)
}

// InitCoinTable func
func InitCoinTable(list []string) *widgets.Table {
	var cryptCoins []Coin
	var c Coin
	for i := 0; i <= len(list)-1; i++ {
		c.Title = list[i]
		c.Value = ""
		cryptCoins = append(cryptCoins, c)
	}
	rows := (&CoinList{Coins: cryptCoins}).coinArrayRows()
	fmt.Printf("coins!\n")
	fmt.Printf("%+v\n", cryptCoins)

	table := widgets.NewTable()

	table.Rows = rows

	return table
}

// MapToRows converts a map of string, string to 2-D rows.
func MapToRows(m map[string]string) [][]string {
	sorted := sortMap(m)
	var headers, values []string
	var rows [][]string
	for colkey, value := range sorted {
		headers = append(headers, colkey)
		values = append(values, value)
	}
	rows = append(rows, headers)
	rows = append(rows, values)
	return rows
}

func sortMap(m map[string]string) map[string]string {
	sorted := make(map[string]string)
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, val := range keys {
		sorted[val] = m[val]
	}

	return m
}

func (c *CoinList) coinArrayRows() [][]string {
	var rows [][]string
	var colHeaders, coinPrices []string

	for i := 0; i <= len(c.Coins)-1; i++ {
		colHeaders = append(colHeaders, (c.Coins)[i].Title)
		coinPrices = append(coinPrices, (c.Coins)[i].Value)
	}
	rows = append(rows, colHeaders)
	rows = append(rows, coinPrices)

	return rows

}

// Map function to act on an array of coins
func Map(vs CoinList, f func(*Coin) string) []string {
	vsm := make([]string, len(vs.Coins))
	for i, v := range vs.Coins {
		vsm[i] = f(&v)
	}
	return vsm
}

// FillCoinsFromMap func
func FillCoinsFromMap(m *map[string]string) *CoinList {
	var key, value string
	var coin Coin
	var coinList []Coin

	for key, value = range *m {
		coin.Title = key
		coin.Value = value
		coinList = append(coinList, coin)
	}
	return &CoinList{Coins: coinList}
}

func testData() CoinList {
	var eth, btc Coin

	eth.Title = "ethereum"
	eth.Value = "1.00"

	btc.Title = "bitcoin"
	btc.Value = "2.990"

	return CoinList{Coins: []Coin{eth, btc}}
}
