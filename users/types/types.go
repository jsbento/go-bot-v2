package types

const (
	NO_NEGATIVE = 1
	BOOST_110   = 2
	BOOST_125   = 3
	BOOST_150   = 4
	BOOST_175   = 5
	BOOST_200   = 6
)

type User struct {
	Username   string     `bson:"username"`
	TokenCount int        `bson:"token_count"`
	PowerUps   []*PowerUp `bson:"power_ups"`
}

type PowerUp struct {
	Value    int  `bson:"value"`
	Active   bool `bson:"active"`
	Modifier int  `bson:"modifier"`
}

type PowerUpSlice []*PowerUp

func SetDefaults() []*PowerUp {
	return PowerUpSlice{
		&PowerUp{Value: 300, Active: false, Modifier: NO_NEGATIVE}, // Limit numbers of rolls active
		&PowerUp{Value: 500, Active: false, Modifier: BOOST_110},
		&PowerUp{Value: 1000, Active: false, Modifier: BOOST_125},
		&PowerUp{Value: 1500, Active: false, Modifier: BOOST_150},
		&PowerUp{Value: 2000, Active: false, Modifier: BOOST_175},
		&PowerUp{Value: 5000, Active: false, Modifier: BOOST_200},
	}
}
