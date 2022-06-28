package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	bU "github.com/jsbento/go-bot-v2/bot/bot_client"
	bot "github.com/jsbento/go-bot-v2/bot/types"
)

func main() {
	bot := bot.Bot{}
	err := bot.New()
	if err == nil {
		defer bot.Close()
		bU.SetDBClient(bot.DbClient)
		fmt.Println("Bot running. Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
	} else {
		fmt.Println("Error creating bot client,", err)
	}
}
