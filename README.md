# nozobot [![GoDoc](https://godoc.org/github.com/VineBalloon/nozobot?status.svg)](https://godoc.org/github.com/VineBalloon/nozobot)

Another Go Discord Bot.

## Dependencies
* [discordgo](https://github.com/bwmarrin/discordgo)
* go 1.12
* Docker (if you want to use that)

## Installation
### with docker
This is probably easier.

0. [Get Docker](https://docs.docker.com/install/)
1. `docker build --tag=nozobot .` (may have to run with sudo)
2. `docker run -e "CARUDO=[YOUR_BOT_TOKEN]" -e "WAIFU=[YOUR_WAIFU2X_TOKEN]" nozobot`

Will probably upload to a public Docker repo at some point.

### with go
If you'd like to modify and develop.

0. [Get Go](https://golang.org/doc/install)
1. `go get -d -v ./...` \*
2. `go install -v ./...`\*
3. `export CARUDO="YOUR_BOT_TOKEN"`
4. `export CARUDO="YOUR_WAIFU2X_TOKEN"`
5. `go run main.go`

\* - [See this for more info](https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations)
