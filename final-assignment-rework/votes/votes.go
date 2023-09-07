package votes

import "time"

type Vote struct {
	VoteID    uint  `json:"vote_id"`
	VoteValue uint  `json:"vote_value"`
	Items     items `json:"items"`
}

type items struct {
	VoterID uint `json:"voter_id"`
	PollID  uint `json:"poll_id"`
}

type Voter struct {
	VoterID     uint        `json:"voter_id"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	VoteHistory []voterPoll `json:"vote_history"`
}

type voterPoll struct {
	PollID   uint      `json:"poll_id"`
	VoteDate time.Time `json:"vote_date"`
}

type pollOption struct {
	PollOptionID   uint   `json:"option_id"`
	PollOptionText string `json:"option_text"`
}

type Poll struct {
	PollID       uint         `json:"poll_id"`
	PollTitle    string       `json:"poll_title"`
	PollQuestion string       `json:"poll_question"`
	PollOptions  []pollOption `json:"poll_options"`
}
