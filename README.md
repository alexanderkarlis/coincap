# CoinCap Price Streamer 💰
A program that runs in a terminal tab and streams cryptocurrency prices from [CoinCap's](https://docs.coincap.io/?version=latest) websocket stream. 

![alt text](https://media.giphy.com/media/MEF2w6McJQ4IX8RDGh/giphy.gif)

## Install
```sh
> go install .
> /path/to/go/bin/coincap bitcoin ethereum litecoin
```

## Build
Clone this repo and build the binary. Command takes and number of args in the form of long coin-names (e.g.- for `Bitcoin` arg should be `bitcoin` and not the short name `btc`)
```sh
> cd coincap/
> go run build
> ./coincap bitcoin ethereum litecoin
> ^C
```
Type or `ctrl+C` or `q` to exit screen.

## Requirements
* Go 1.14.x

## Known issues
- [ ] -- coin table columns are shifting around making it hard to track which coins are actually getting updated. 
- [ ] -- sometimes exiting the screen isn't always responsive. In this case, kill the process. 

