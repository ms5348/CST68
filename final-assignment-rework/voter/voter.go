package voter

import (
	"time"
)

type VoterPoll struct {
	PollID   uint      `json:"poll_id"`
	VoteDate time.Time `json:"vote_date"`
}

type Voter struct {
	VoterID     uint        `json:"voter_id"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	VoteHistory []VoterPoll `json:"vote_history"`
}
