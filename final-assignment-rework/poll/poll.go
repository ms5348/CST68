package poll

type PollOption struct {
	PollOptionID   uint   `json:"option_id"`
	PollOptionText string `json:"option_text"`
}

type Poll struct {
	PollID       uint         `json:"poll_id"`
	PollTitle    string       `json:"poll_title"`
	PollQuestion string       `json:"poll_question"`
	PollOptions  []PollOption `json:"poll_options"`
}
