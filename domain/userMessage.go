package domain

// UserMessage messageType 201 user created
// messageType 200 user updated
// messageType 204 user deleted
type UserMessage struct {
	User *User `form:"user" json:"user"`
	MessageType int `form:"messageType" json:"messageType"`
	ResourceType string `form:"resourceType" json:"resourceType"`
}
