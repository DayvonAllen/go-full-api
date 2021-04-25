package domain

type UserMessage struct {
	User User `form:"user" json:"user"`
}
