package bot_client

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	dbT "github.com/jsbento/go-bot-v2/db/types"
	uT "github.com/jsbento/go-bot-v2/users/types"
	uU "github.com/jsbento/go-bot-v2/users/user_utils"
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
			PowerUps:   uT.SetDefaults(),
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
		tokens = uU.ApplyPowerUps(user.PowerUps, tokens)
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

	// Show current balance and powerups owned
	if m.Content == "!powerups" {
		_, err := s.ChannelMessageSend(m.ChannelID,
			"Powerups:\n"+
				"pId 1: No-Negatives - Negatives count as 0's. Cost: 300 tokens\n"+
				"pId 2: 110% Boost. Cost: 500 tokens\n"+
				"pId 3: 125% Boost. Cost: 1000 tokens\n"+
				"pId 4: 150% Boost. Cost: 1500 tokens\n"+
				"pId 5: 175% Boost. Cost: 2000 tokens\n"+
				"pId 6: 200% Boost. Cost: 3000 tokens\n"+
				"To purchase a powerup, type !powerup buy <pId>\n")
		if err != nil {
			fmt.Println("Error sending message,", err)
			return
		}
	}

	if strings.Contains(m.Content, "!powerup buy ") {
		pId, err := strconv.Atoi(strings.Split(m.Content, " ")[2])
		if err != nil {
			fmt.Println("Error converting pId,", err)
			return
		}
		user, err := dbClient.GetUser(m.Author.Username)
		if err != nil {
			if user.Username == "" {
				_, err = s.ChannelMessageSend(m.ChannelID, "You don't have a saved user!")
			}
			fmt.Println("Error getting user,", err)
			return
		}
		result, err := dbClient.PurchasePowerUp(user.Username, pId)
		if err != nil {
			if result == nil {
				fmt.Println("Error purchasing powerup,", err)
				return
			} else {
				_, err = s.ChannelMessageSend(m.ChannelID, err.Error())
				if err != nil {
					fmt.Println("Error sending message,", err)
					return
				}
				return
			}
		}
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully purchased powerup %d", pId))
		if err != nil {
			fmt.Println("Error sending message,", err)
			return
		}
	}
}
