package votes

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

type Vote struct {
	VoteID    uint `json:"vote_id"`
	VoterID   uint `json:"voter_id"`
	PollID    uint `json:"poll_id"`
	VoteValue uint `json:"vote_value"`
}

type VoteCache struct {
	cache
}

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "vote:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

func New() (*VoteCache, error) {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	return NewWithCacheInstance(redisUrl)
}

func NewWithCacheInstance(location string) (*VoteCache, error) {

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

	return &VoteCache{
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

func (v *VoteCache) getVoteFromRedis(key string, vote *Vote) error {
	voteObject, err := v.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(voteObject.([]byte), vote)
	if err != nil {
		return err
	}

	return nil
}

func (v *VoteCache) AddVote(vote Vote) error {
	redisKey := redisKeyFromId(vote.VoteID)
	var existingVote Vote
	if err := v.getVoteFromRedis(redisKey, &existingVote); err == nil {
		return errors.New("vote already exists")
	}

	if _, err := v.jsonHelper.JSONSet(redisKey, ".", vote); err != nil {
		return err
	}

	return nil
}

func (v *VoteCache) GetVote(id uint) (Vote, error) {
	var vote Vote
	pattern := redisKeyFromId(id)
	err := v.getVoteFromRedis(pattern, &vote)
	if err != nil {
		return Vote{}, err
	}

	return vote, nil
}

func (v *VoteCache) GetAllVotes() ([]Vote, error) {
	var voteList []Vote
	var vote Vote

	pattern := RedisKeyPrefix + "*"
	ks, _ := v.cacheClient.Keys(v.context, pattern).Result()
	for _, key := range ks {
		err := v.getVoteFromRedis(key, &vote)
		if err != nil {
			return nil, err
		}
		voteList = append(voteList, vote)
	}

	return voteList, nil
}

func (v *Vote) ToJson() string {
	b, _ := json.Marshal(v)
	return string(b)
}

func NewSampleVote() *Vote {
	return &Vote{
		VoteID:    1,
		PollID:    1,
		VoterID:   1,
		VoteValue: 1,
	}
}
