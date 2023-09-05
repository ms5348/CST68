package poll

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

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

type PollCache struct {
	cache
}

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "poll:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

func New() (*PollCache, error) {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	return NewWithCacheInstance(redisUrl)
}

func NewWithCacheInstance(location string) (*PollCache, error) {

	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	ctx := context.Background()

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err
	}

	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	return &PollCache{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}, nil
}

func redisKeyFromId(id uint) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

func (p *PollCache) getPollFromRedis(key string, poll *Poll) error {
	pollObject, err := p.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(pollObject.([]byte), poll)
	if err != nil {
		return err
	}

	return nil
}

func (p *PollCache) AddPoll(poll Poll) error {
	redisKey := redisKeyFromId(poll.PollID)
	var existingPoll Poll
	if err := p.getPollFromRedis(redisKey, &existingPoll); err == nil {
		return errors.New("poll already exists")
	}

	if _, err := p.jsonHelper.JSONSet(redisKey, ".", poll); err != nil {
		return err
	}

	return nil
}

/*
func (v *VoterCache) AddPoll(voter Voter) error {
	redisKey := redisKeyFromId(voter.VoterID)
	var existingVoter Voter
	if err := v.getVoterFromRedis(redisKey, &existingVoter); err != nil {
		return errors.New("voter does not exist")
	}

	poll := voter.VoteHistory[pollID]
	if poll ==  {
		return poll, err
	}

	if _, err := v.jsonHelper.JSONSet(redisKey, ".", voter); err != nil {
		return err
	}

	return nil
}*/

func (p *PollCache) GetPoll(id uint) (Poll, error) {
	var poll Poll
	pattern := redisKeyFromId(id)
	err := p.getPollFromRedis(pattern, &poll)
	if err != nil {
		return Poll{}, err
	}

	return poll, nil
}

func (p *PollCache) GetAllPolls() ([]Poll, error) {
	var pollList []Poll
	var poll Poll

	pattern := RedisKeyPrefix + "*"
	ks, _ := p.cacheClient.Keys(p.context, pattern).Result()
	for _, key := range ks {
		err := p.getPollFromRedis(key, &poll)
		if err != nil {
			return nil, err
		}
		pollList = append(pollList, poll)
	}

	return pollList, nil
}

func (p *PollCache) GetPollOptions(id uint) ([]pollOption, error) {
	var poll Poll
	var pollOptions []pollOption
	pattern := redisKeyFromId(id)
	err := p.getPollFromRedis(pattern, &poll)
	if err != nil {
		return pollOptions, err
	}

	return poll.PollOptions, nil
}

/*
func (v *VoterCache) GetPoll(voterID uint, pollID uint) (voterPoll, error) {
	var voter Voter
	var poll voterPoll
	pattern := redisKeyFromId(voterID)
	err := v.getVoterFromRedis(pattern, &voter)
	if err != nil {
		return poll, err
	}

	// I do not believe the following if statement is functioning as intended
	if pollID > uint(len(voter.VoteHistory)) {
		return poll, errors.New(" redis: nil")
	}
	poll = voter.VoteHistory[pollID-1]
	return poll, nil
}*/

func (p *Poll) ToJson() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func NewSamplePoll() *Poll {
	return &Poll{
		PollID:       1,
		PollTitle:    "Favorite Pet",
		PollQuestion: "What type of pet do you like best?",
		PollOptions: []pollOption{
			{PollOptionID: 1, PollOptionText: "Dog"},
			{PollOptionID: 2, PollOptionText: "Cat"},
			{PollOptionID: 3, PollOptionText: "Fish"},
			{PollOptionID: 4, PollOptionText: "Bird"},
			{PollOptionID: 5, PollOptionText: "NONE"},
		},
	}
}
