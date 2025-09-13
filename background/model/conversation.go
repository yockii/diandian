package model

type Conversation struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	ID             uint   `json:"id"`
	ConversationID uint   `json:"conversation_id"`
	Role           string `json:"role"` // user, assistant
	Content        string `json:"content"`
}
