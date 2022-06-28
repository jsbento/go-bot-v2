package types

import (
	"errors"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	bC "github.com/jsbento/go-bot-v2/bot/bot_client"
	dbT "github.com/jsbento/go-bot-v2/db/types"
)

type Bot struct {
	DiscordClient *discordgo.Session
	DbClient      *dbT.DBClient
}

func (b *Bot) New() error {
	err := godotenv.Load(".env")
	if err != nil {
		return errors.New("could not load .env file")
	}
	botToken := os.Getenv("BOT_TOKEN")

	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return fmt.Errorf("error creating discord session, %s", err)
	}

	err = dg.Open()
	if err != nil {
		return fmt.Errorf("error opening discord connection, %s", err)
	}

	dg.AddHandler(bC.MessageCreate)
	dg.Identify.Intents = discordgo.IntentGuildMessages

	dbClient, err := dbT.New()
	if err != nil {
		dg.Close()
		return fmt.Errorf("error creating db client, %s", err)
	}

	b.DiscordClient = dg
	b.DbClient = dbClient
	return nil
}

func (b Bot) Close() {
	fmt.Println("Closing bot connections...")
	b.DbClient.Close()
	b.DiscordClient.Close()
}
