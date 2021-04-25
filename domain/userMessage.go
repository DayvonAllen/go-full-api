package domain

// UserMessage messageType 201 user created
// messageType 200 user updated
type UserMessage struct {
	User User `form:"user" json:"user"`
	MessageType int `form:"messageType" json:"messageType"`
}
