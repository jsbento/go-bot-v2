package users

import (
	uT "github.com/jsbento/go-bot-v2/users/types"
)

const (
	NO_NEGATIVE = 1
	BOOST_110   = 2
	BOOST_125   = 3
	BOOST_150   = 4
	BOOST_175   = 5
	BOOST_200   = 6
)

func ApplyPowerUps(powerups []*uT.PowerUp, tokens int) (newTokens int) {
	newTokens = tokens
	boostTotal := 1.0
	for _, powerup := range powerups {
		if powerup.Active {
			switch powerup.Modifier {
			case NO_NEGATIVE:
				if tokens < 0 {
					powerup.Uses -= 1
					if powerup.Uses <= 0 {
						powerup.Active = false
					}
					return 0
				}
			case BOOST_110:
				boostTotal += 0.1
			case BOOST_125:
				boostTotal += 0.25
			case BOOST_150:
				boostTotal += 0.5
			case BOOST_175:
				boostTotal += 0.75
			case BOOST_200:
				boostTotal += 1.0
			default:
			}
		}
	}
	return int(float64(tokens) * boostTotal)
}

func GetActivePowerUps(poweups []*uT.PowerUp) (activePowerups []*uT.PowerUp) {
	for _, powerup := range poweups {
		if powerup.Active {
			activePowerups = append(activePowerups, powerup)
		}
	}
	return
}
