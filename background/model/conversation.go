package model

type Conversation struct {
	Base
	Name string `json:"name" gorm:"size:200"`
}

type Message struct {
	Base
	ConversationID uint64 `json:"conversation_id" gorm:"index"`
	Role           string `json:"role" gorm:"size:50"` // user, assistant
	Content        string `json:"content"`
}
