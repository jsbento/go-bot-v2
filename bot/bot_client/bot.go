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

const ( // dbClient flags
	DB_ERROR        = -1
	USER_OK         = 0
	USER_EXISTS     = 1
	USER_NOT_EXISTS = 2
	LOW_TOKENS      = 3
	POWER_ACTIVE    = 4
)

const ( // Calling functions flags, affects messages
	CREATE    = 0
	GET       = 1
	DELETE    = 2
	ROLL      = 3
	POWERS    = 4
	POWER_BUY = 5
	ROLL_ALT  = 6
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

	// Have each dbClient function return (flag int, err error) -> pass to helper function to send messages based on flag
	// Should help cleanup some repeated code
	if m.Content == "!create_me" {
		user := uT.User{
			Username:   m.Author.Username,
			TokenCount: 0,
			PowerUps:   uT.SetDefaults(),
		}

		flag, err := dbClient.PostUser(&user)
		HandleMessage(s, m, flag, CREATE, 0, err, user)
		return
	}

	if m.Content == "!delete_me" {
		flag, err := dbClient.DeleteUser(m.Author.Username)
		HandleMessage(s, m, flag, DELETE, 0, err, uT.User{Username: m.Author.Username})
		return
	}

	if m.Content == "!get_me" {
		user, flag, err := dbClient.GetUser(m.Author.Username)
		HandleMessage(s, m, flag, GET, 0, err, user)
		return
	}

	if m.Content == "!roll" {
		user, flag, err := dbClient.GetUser(m.Author.Username)
		HandleMessage(s, m, flag, ROLL_ALT, 0, err, user)
		seed := rand.NewSource(time.Now().UnixNano())
		r := rand.New(seed)
		tokens := r.Intn(250) - 75
		tokens = uU.ApplyPowerUps(user.PowerUps, tokens)
		flag, err = dbClient.UpdateUser(user.Username, user.PowerUps, tokens)
		HandleMessage(s, m, flag, ROLL, tokens, err, user)
		return
	}

	if m.Content == "!powerups" {
		user, flag, err := dbClient.GetUser(m.Author.Username)
		HandleMessage(s, m, flag, POWERS, 0, err, user)
		return
	}

	if strings.Contains(m.Content, "!powerup buy ") {
		pId, err := strconv.Atoi(strings.Split(m.Content, " ")[2])
		if err != nil {
			fmt.Println("Error converting pId,", err)
			return
		}
		user, flag, err := dbClient.GetUser(m.Author.Username)
		HandleMessage(s, m, flag, POWER_BUY, pId, err, user)
		result, flag, err := dbClient.PurchasePowerUp(user.Username, pId)
		HandleMessage(s, m, flag, POWER_BUY, pId, err, *result)
		return
	}
}

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate, flag, caller, value int, err error, user ...uT.User) {
	var message string
	switch flag {
	case DB_ERROR:
		fmt.Println("Database error: ", err)
		return
	case USER_EXISTS:
		message = err.Error()
	case USER_NOT_EXISTS:
		message = err.Error()
	case LOW_TOKENS:
		message = err.Error()
	case POWER_ACTIVE:
		message = err.Error()
	case USER_OK:
		switch caller {
		case CREATE:
			message = fmt.Sprintf("User %s created successfully with %d tokens", user[0].Username, user[0].TokenCount)
		case GET:
			message = fmt.Sprintf("You have %d tokens", user[0].TokenCount)
		case DELETE:
			message = fmt.Sprintf("User %s deleted", user[0].Username)
		case ROLL:
			message = fmt.Sprintf("You rolled %d tokens", value)
		case POWERS:
			var activeIds []string
			active := uU.GetActivePowerUps(user[0].PowerUps)
			for _, v := range active {
				activeIds = append(activeIds, strconv.Itoa(v.Modifier))
			}
			powerMsg := "Powerups:\n" +
				"pId 1: No-Negatives - Negatives count as 0's. Cost: 300 tokens\n" +
				"pId 2: 110% Boost. Cost: 500 tokens\n" +
				"pId 3: 125% Boost. Cost: 1000 tokens\n" +
				"pId 4: 150% Boost. Cost: 1500 tokens\n" +
				"pId 5: 175% Boost. Cost: 2000 tokens\n" +
				"pId 6: 200% Boost. Cost: 3000 tokens\n" +
				"To purchase a powerup, type !powerup buy <pId>"
			message = fmt.Sprintf("Current balance: %d, Active: %s\n%s",
				user[0].TokenCount,
				strings.Join(activeIds, ", "),
				powerMsg)
		case POWER_BUY:
			message = fmt.Sprintf("Successfully purchased powerup %d", value)
		}
	default:
		return
	}
	_, err = s.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		fmt.Println("Error sending message,", err)
		return
	}
}
