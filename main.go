package main

import (
	"os"

	"github.com/alexanderkarlis/coincap/rabbit"
	"github.com/alexanderkarlis/coincap/stream"
)

func main() {
	coins := os.Args[1:]
	go stream.Stream(coins)
	rabbit.Consume()
}
