package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"final-assignment/voter"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

type VoterAPI struct {
	cache
}

func NewVoterAPI(location string) (*VoterAPI, error) {
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

	return &VoterAPI{
		cache: cache{
			client:  client,
			helper:  jsonHelper,
			context: ctx,
		},
	}, nil
}

func (v *VoterAPI) AddVoter(c *gin.Context) {
	var newVoter voter.Voter

	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No voter ID provided"})
		return
	}

	cacheKey := "voter:" + id

	var existingVoter voter.Voter
	if err, _ := v.getItemFromRedis(cacheKey, &existingVoter); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Voter ID already exists"})
		return
	}

	if _, err := v.helper.JSONSet(cacheKey, ".", newVoter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new Voter"})
		return
	}

	c.JSON(http.StatusOK, newVoter)
}

func (v *VoterAPI) AddPoll(c *gin.Context) {
	var voterItem voter.Voter

	if err := c.ShouldBindJSON(&voterItem); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No voter ID provided"})
		return
	}
	cacheKey := "voter:" + id

	var existingVoter voter.Voter
	err, existingVoterHistory := v.getItemFromRedis(cacheKey, &existingVoter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find voter in cache with id=" + cacheKey})
		return
	}

	pollIDS := c.Param("pollid")
	if pollIDS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll ID provided"})
		return
	}

	pollID64, err := strconv.ParseUint(pollIDS, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to iunt64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollID := uint(pollID64)

	for _, voterPoll := range existingVoterHistory {
		if voterPoll.PollID == pollID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID already exists"})
			return
		}
	}

	if _, err := v.helper.JSONSet(cacheKey, ".", voterItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new VoterHistory"})
		return
	}

	c.JSON(http.StatusOK, voterItem)
}

func (v *VoterAPI) GetVoter(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Voter ID provided"})
		return
	}

	cacheKey := "voter:" + id
	voterBytes, err := v.helper.JSONGet(cacheKey, ".")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find voter in cache with id=" + cacheKey})
		return
	}

	var voter voter.Voter
	err = json.Unmarshal(voterBytes.([]byte), &voter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cached data seems to be wrong type"})
		return
	}

	c.JSON(http.StatusOK, voter)
}

func (v *VoterAPI) GetVoterList(c *gin.Context) {
	var voterList []voter.Voter
	var voterItem voter.Voter

	pattern := "voter:*"
	ks, _ := v.client.Keys(v.context, pattern).Result()
	for _, key := range ks {
		err, _ := v.getItemFromRedis(key, &voterItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find voter in cache with id=" + key})
			return
		}
		voterList = append(voterList, voterItem)
	}

	c.JSON(http.StatusOK, voterList)
}

func (v *VoterAPI) GetPolls(c *gin.Context) {
	id := c.Param("id")
	var thisVoter voter.Voter
	pattern := "voter:" + id
	err, _ := v.getItemFromRedis(pattern, &thisVoter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find voter in cache with id=" + pattern})
		return
	}

	c.JSON(http.StatusOK, thisVoter.VoteHistory)
}

func (v *VoterAPI) GetPoll(c *gin.Context) {
	id := c.Param("id")
	pollIDS := c.Param("pollid")
	if pollIDS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll ID provided"})
		return
	}
	pollID64, err := strconv.ParseInt(pollIDS, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	cacheKey := "voter:" + id

	var existingVoter voter.Voter
	err, existingVoterHistory := v.getItemFromRedis(cacheKey, &existingVoter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find voter in cache with id=" + cacheKey})
		return
	}

	pollID := uint(pollID64)

	for _, voterPoll := range existingVoterHistory {
		if voterPoll.PollID == pollID {
			c.JSON(http.StatusOK, voterPoll)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll in cache with id=" + pollIDS})
}

func (v *VoterAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"status":             "ok",
			"version":            "1.0.0",
			"uptime":             100,
			"users_processed":    1000,
			"errors_encountered": 10,
		})
}

func (v *VoterAPI) getItemFromRedis(key string, voter *voter.Voter) (error, []voter.VoterPoll) {
	itemObject, err := v.helper.JSONGet(key, ".")
	if err != nil {
		return err, voter.VoteHistory
	}

	err = json.Unmarshal(itemObject.([]byte), voter)
	if err != nil {
		return err, voter.VoteHistory
	}

	return nil, voter.VoteHistory
}
