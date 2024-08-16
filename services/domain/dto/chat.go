package dto

type ChatBotSubmitDto struct {
	ConversationId string
	Question       string
	UserId         string
	RequestId      string
}

type FeedbackAnswerDto struct {
	MessageId      string
	Feedback       string
	UserId         string
	ConversationId string
}

type DetailAnswerDto struct {
	ConversationId string
	UserId         string
}

type ListConversationDto struct {
	Limit  uint32
	Page   uint32
	UserId string
}
