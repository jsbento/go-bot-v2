package types

type User struct {
	Username   string `bson:"username"`
	TokenCount int    `bson:"token_count"`
}
