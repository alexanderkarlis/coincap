# CoinCap Price Streamer 💰
A program that runs in a terminal tab and streams cryptocurrency prices from [CoinCap's](https://docs.coincap.io/?version=latest) websocket stream. 

![alt text](https://media.giphy.com/media/Ln41AuY0aNBlLeo59I/giphy.gif)

#### *Notes*:
 Having RabbitMQ handle the message busing is overkill, yes. I was curious about it, and I 
 was thinking of extending this project later. So, yeah, two birds.

Starts a channel in RabbitMQ under the `coin-prices` so in theory extra applications could query this channel, thus building ontop. 

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
* RabbitMQ running on `localhost:5672`. (**TODO**: make more dynamic via env, settings, etc..)

## Known issues
- [] -- coin table columns are shifting around making it hard to track which coins are actually getting updated. 
- [] -- sometimes exiting the screen isn't always responsive. In this case, kill the process. 

