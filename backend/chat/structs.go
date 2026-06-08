package chat

type getMessagesInput struct {
	ID       uint `path:"id" doc:"Event ID" example:"42"`
	BeforeID uint `query:"before_id" default:"0" doc:"Load messages before this message ID" example:"120"`
}

type createMessageInput struct {
	Content string `json:"content" doc:"Message content" example:"hello chat"`
}

type MessageListDTO struct {
	Data []MessageDTO `json:"data" doc:"Chat messages"`
}

type messagesOutput struct {
	Body MessageListDTO
}
