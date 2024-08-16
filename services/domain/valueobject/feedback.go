package valueobject

type FeedbackEnum string

const (
	LIKE       FeedbackEnum = "LIKE"
	DISLIKE    FeedbackEnum = "DISLIKE"
	NOT_DEFINE FeedbackEnum = "NOT_DEFINE"
)

type Feedback struct {
	value FeedbackEnum
}

func NewFeedBack(s FeedbackEnum) Feedback {
	return Feedback{value: s}
}

func NewFeedBackFromParam(s string) Feedback {
	var valueEnum FeedbackEnum
	switch s {
	case "LIKE":
		valueEnum = LIKE
	case "DISLIKE":
		valueEnum = DISLIKE
	default:
		valueEnum = NOT_DEFINE
	}

	return Feedback{value: valueEnum}
}

func (i Feedback) String() string {
	result := ""
	switch i.value {
	case LIKE:
		result = "like"
	case DISLIKE:
		result = "dislike"
	}

	return result
}
