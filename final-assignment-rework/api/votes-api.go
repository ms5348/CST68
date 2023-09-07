package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"final-assignment/votes"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/nitishm/go-rejson/v4"
)

type cache struct {
	client  *redis.Client
	helper  *rejson.Handler
	context context.Context
}

type VotesAPI struct {
	cache
	voterAPIURL string
	pollAPIURL  string
	apiClient   *resty.Client
}

func NewVotesAPI(location string, voterAPIurl string, pollAPIURL string) (*VotesAPI, error) {
	apiClient := resty.New()

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

	return &VotesAPI{
		cache: cache{
			client:  client,
			helper:  jsonHelper,
			context: ctx,
		},
		voterAPIURL: voterAPIurl,
		pollAPIURL:  pollAPIURL,
		apiClient:   apiClient,
	}, nil
}

func (v *VotesAPI) AddVote(c *gin.Context) {
	var newVote votes.Vote

	if err := c.ShouldBindJSON(&newVote); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No vote ID provided"})
		return
	}

	cacheKey := "vote:" + id

	var existingVote votes.Vote
	if err := v.getItemFromRedis(cacheKey, &existingVote); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vote ID already exists"})
		return
	}

	if _, err := v.helper.JSONSet(cacheKey, ".", newVote); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new Vote"})
		return
	}

	c.JSON(http.StatusOK, newVote)
}

func (v *VotesAPI) GetVotes(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No vote ID provided"})
		return
	}

	cacheKey := "vote:" + id
	voteBytes, err := v.helper.JSONGet(cacheKey, ".")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find vote in cache with id=" + cacheKey})
		return
	}

	var vote votes.Vote
	err = json.Unmarshal(voteBytes.([]byte), &vote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cached data seems to be wrong type"})
		return
	}

	c.JSON(http.StatusOK, vote)
}

func (v *VotesAPI) GetVotesList(c *gin.Context) {
	var voteList []votes.Vote
	var voteItem votes.Vote

	pattern := "vote:*"
	ks, _ := v.client.Keys(v.context, pattern).Result()
	for _, key := range ks {
		err := v.getItemFromRedis(key, &voteItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find reading list in cache with id=" + key})
			return
		}
		voteList = append(voteList, voteItem)
	}

	c.JSON(http.StatusOK, voteList)
}

func (v *VotesAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"status":             "ok",
			"version":            "2.0.0",
			"uptime":             200,
			"users_processed":    2000,
			"errors_encountered": 20,
		})
}

func (v *VotesAPI) GetItemFromVote(c *gin.Context) {
	voteId := c.Param("id")
	if voteId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No vote ID provided"})
		return
	}

	idxKey := c.Param("idx")
	if idxKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No item index provided"})
		return
	}

	cacheKey := "vote:" + voteId
	var vote votes.Vote
	err := v.getItemFromRedis(cacheKey, &vote)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find vote in cache with id=" + cacheKey})
		return
	}

	if idxKey == "voter_id" {
		voterLocation := vote.Items.VoterID

		voterURL := v.voterAPIURL + "/" + strconv.FormatUint(uint64(voterLocation), 10)
		var voter votes.Voter

		_, err = v.apiClient.R().SetResult(&voter).Get(voterURL)
		if err != nil {
			emsg := "Could not get voter from API: (" + voterURL + ")" + err.Error()
			c.JSON(http.StatusNotFound, gin.H{"error": emsg})
			return
		}

		c.JSON(http.StatusOK, voter)
		return
	}

	if idxKey == "poll_id" {
		pollLocation := vote.Items.PollID

		pollURL := v.pollAPIURL + "/" + strconv.FormatUint(uint64(pollLocation), 10)
		var poll votes.Poll

		_, err = v.apiClient.R().SetResult(&poll).Get(pollURL)
		if err != nil {
			emsg := "Could not get poll from API: (" + pollURL + ")" + err.Error()
			c.JSON(http.StatusNotFound, gin.H{"error": emsg})
			return
		}

		c.JSON(http.StatusOK, poll)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect index provided"})
}

func (v *VotesAPI) getItemFromRedis(key string, vote *votes.Vote) error {
	itemObject, err := v.helper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(itemObject.([]byte), vote)
	if err != nil {
		return err
	}

	return nil
}
