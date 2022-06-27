package main

import (
	"github.com/jsbento/go-bot-v2/bot"
	"github.com/jsbento/go-bot-v2/mongo"
)

func main() {
	bot.New(mongo.New())
}
