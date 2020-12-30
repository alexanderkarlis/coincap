# CoinCap Price Streamer 💰
A program that runs in a terminal tab and streams cryptocurrency prices from [CoinCap's](https://docs.coincap.io/?version=latest) websocket stream. 

![alt text](https://media.giphy.com/media/MEF2w6McJQ4IX8RDGh/giphy.gif)

## Install
```sh
> go install .
> /path/to/go/bin/coincap bitcoin ethereum litecoin
```

## Build & Use
Clone this repo and build the binary. Command takes and number of args in the form of long coin-names (e.g.- for `Bitcoin` arg should be `bitcoin` and not the short name `btc`)
```sh
> cd coincap/
> go run build
> ./coincap bitcoin ethereum litecoin
> ^C
```
This creates a list of tabs of the passed in coins for you to navigate between with `h` and `l` keys
(vim left/right cursor move).
Type or `ctrl+C` or `q` to exit screen.

## Requirements
* Go 1.14.x

## Known issues
- [ ] -- plot of coin prices needs to be ranged on the Y axis better
- [ ] -- move the incoming plots to the ~90 most recent 
