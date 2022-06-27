package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const DogAPIURL = "https://dog.ceo/api/"

func New() {
	botToken := os.Getenv("BOT_TOKEN")

	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	fmt.Println("Bot running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	if m.Content == "!dog" {
		res, err := http.Get(DogAPIURL + "breeds/image/random")
		if err != nil {
			fmt.Println("Error getting dog,", err)
			return
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			fmt.Println("Error getting dog,", res.StatusCode)
			return
		} else {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println("Error reading dog,", err)
				return
			}
			var response Response
			err = json.Unmarshal(body, &response)
			if err != nil {
				fmt.Println("Error unmarshalling dog,", err)
				return
			}
			_, err = s.ChannelMessageSend(m.ChannelID, response.Message)
			if err != nil {
				fmt.Println("Error sending message,", err)
				return
			}
		}
	}
}
