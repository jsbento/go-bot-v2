package bot_client

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	dbT "github.com/jsbento/go-bot-v2/db/types"
	uT "github.com/jsbento/go-bot-v2/users/types"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var dbClient *dbT.DBClient

func SetDBClient(client *dbT.DBClient) {
	dbClient = client
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!me" {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hello, you are: %s", m.Author.Username))
		if err != nil {
			fmt.Println("Error sending message,", err)
			return
		}
	}

	if m.Content == "!create_me" {
		user := uT.User{
			Username:   m.Author.Username,
			TokenCount: 0,
		}

		flag, err := dbClient.PostUser(&user)
		if err != nil {
			if flag == 1 {
				_, err = s.ChannelMessageSend(m.ChannelID, "You already have a saved user!")
			}
			fmt.Println("Error creating user,", err)
			return
		}
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s created with %d tokens", user.Username, user.TokenCount))
		if err != nil {
			fmt.Println("Error sending message,", err)
			return
		}
	}

	if m.Content == "!delete_me" {
		flag, err := dbClient.DeleteUser(m.Author.Username)
		if err != nil {
			if flag == 1 {
				_, err = s.ChannelMessageSend(m.ChannelID, "You don't have a saved user!")
			}
			fmt.Println("Error deleting user,", err)
			return
		}
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s deleted", m.Author.Username))
		if err != nil {
			fmt.Println("Error sending message,", err)
			return
		}
	}

	if m.Content == "!get_me" {
		user, err := dbClient.GetUser(m.Author.Username)
		if err != nil {
			if user.Username == "" {
				_, err = s.ChannelMessageSend(m.ChannelID, "You don't have a saved user!")
			}
			fmt.Println("Error getting user,", err)
			return
		}
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You have %d tokens", user.TokenCount))
		if err != nil {
			fmt.Println("Error sending message,", err)
			return
		}
	}

	// Implement power-ups
	// 1. No negatives -> if negative is rolled, it doesn't count
	// 2. Multiplier -> Multiply roll by some boost (look at stabable boosts)

	if m.Content == "!roll" {
		user, err := dbClient.GetUser(m.Author.Username)
		if err != nil {
			if user.Username == "" {
				_, err = s.ChannelMessageSend(m.ChannelID, "You don't have a saved user!")
			}
			fmt.Println("Error getting user,", err)
			return
		}
		seed := rand.NewSource(time.Now().UnixNano())
		r := rand.New(seed)
		tokens := r.Intn(200) - 75
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You rolled %d tokens", tokens))
		if err != nil {
			fmt.Println("Error sending message,", err)
			return
		}
		_, err = dbClient.AddTokens(user.Username, tokens)
		if err != nil {
			fmt.Println("Error updating user,", err)
			return
		}
	}
}
