package voter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

type voterPoll struct {
	PollID   uint      `json:"poll_id"`
	VoteDate time.Time `json:"vote_date"`
}

type Voter struct {
	VoterID     uint        `json:"voter_id"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	VoteHistory []voterPoll `json:"vote_history"`
}
type VoterCache struct {
	cache
}

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "voter:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

func New() (*VoterCache, error) {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	return NewWithCacheInstance(redisUrl)
}

func NewWithCacheInstance(location string) (*VoterCache, error) {

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

	return &VoterCache{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}, nil
}

/*func isRedisNilError(err error) bool {
	return errors.Is(err, redis.Nil) || err.Error() == RedisNilError
}*/

func redisKeyFromId(id uint) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

func (v *VoterCache) getVoterFromRedis(key string, voter *Voter) error {
	voterObject, err := v.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(voterObject.([]byte), voter)
	if err != nil {
		return err
	}

	return nil
}

func (v *VoterCache) AddVoter(voter Voter) error {
	redisKey := redisKeyFromId(voter.VoterID)
	var existingVoter Voter
	if err := v.getVoterFromRedis(redisKey, &existingVoter); err == nil {
		return errors.New("voter already exists")
	}

	if _, err := v.jsonHelper.JSONSet(redisKey, ".", voter); err != nil {
		return err
	}

	return nil
}

func (v *VoterCache) AddPoll(voter Voter) error {
	redisKey := redisKeyFromId(voter.VoterID)
	var existingVoter Voter
	if err := v.getVoterFromRedis(redisKey, &existingVoter); err != nil {
		return errors.New("voter does not exist")
	}

	/*poll := voter.VoteHistory[pollID]
	if poll ==  {
		return poll, err
	}*/

	if _, err := v.jsonHelper.JSONSet(redisKey, ".", voter); err != nil {
		return err
	}

	return nil
}

func (v *VoterCache) GetVoter(id uint) (Voter, error) {
	var voter Voter
	pattern := redisKeyFromId(id)
	err := v.getVoterFromRedis(pattern, &voter)
	if err != nil {
		return Voter{}, err
	}

	return voter, nil
}

func (v *VoterCache) GetAllVoters() ([]Voter, error) {
	var voterList []Voter
	var voter Voter

	pattern := RedisKeyPrefix + "*"
	ks, _ := v.cacheClient.Keys(v.context, pattern).Result()
	for _, key := range ks {
		err := v.getVoterFromRedis(key, &voter)
		if err != nil {
			return nil, err
		}
		voterList = append(voterList, voter)
	}

	return voterList, nil
}

func (v *VoterCache) GetPolls(id uint) ([]voterPoll, error) {
	var voter Voter
	var polls []voterPoll
	pattern := redisKeyFromId(id)
	err := v.getVoterFromRedis(pattern, &voter)
	if err != nil {
		return polls, err
	}

	return voter.VoteHistory, nil
}

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
}

func (v *Voter) ToJson() string {
	b, _ := json.Marshal(v)
	return string(b)
}

func NewSampleVoter() *Voter {
	return &Voter{
		VoterID:   1,
		FirstName: "John",
		LastName:  "Doe",
		VoteHistory: []voterPoll{
			{PollID: 1, VoteDate: time.Now()},
		},
	}
}
