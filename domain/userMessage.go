package domain

// Message messageType 201 user created
// messageType 200 user updated
// messageType 204 user deleted
type Message struct {
	User User `form:"User" json:"User"`
	Event        Event  `form:"Event" json:"Event"`
	MessageType int `form:"messageType" json:"messageType"`
	ResourceType string `form:"resourceType" json:"resourceType"`
}
