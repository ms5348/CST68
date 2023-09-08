package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"final-assignment/poll"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

type PollAPI struct {
	cache
}

func NewPollAPI(location string) (*PollAPI, error) {
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

	return &PollAPI{
		cache: cache{
			client:  client,
			helper:  jsonHelper,
			context: ctx,
		},
	}, nil
}

func (p *PollAPI) AddPoll(c *gin.Context) {
	var newPoll poll.Poll

	if err := c.ShouldBindJSON(&newPoll); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll ID provided"})
		return
	}

	cacheKey := "poll:" + id

	var existingPoll poll.Poll
	if err, _ := p.getItemFromRedis(cacheKey, &existingPoll); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID already exists"})
		return
	}

	if _, err := p.helper.JSONSet(cacheKey, ".", newPoll); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new Poll"})
		return
	}

	c.JSON(http.StatusOK, newPoll)
}

func (p *PollAPI) AddPollOption(c *gin.Context) {
	var pollItem poll.Poll

	if err := c.ShouldBindJSON(&pollItem); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll ID provided"})
		return
	}
	cacheKey := "poll:" + id

	var existingPoll poll.Poll
	err, existingPollOptions := p.getItemFromRedis(cacheKey, &existingPoll)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll in cache with id=" + cacheKey})
		return
	}

	pollIDS := c.Param("pollid")
	if pollIDS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll option ID provided"})
		return
	}

	pollID64, err := strconv.ParseUint(pollIDS, 10, 32)
	if err != nil {
		log.Println("Error converting poll option id to iunt64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollID := uint(pollID64)

	for _, pollOption := range existingPollOptions {
		if pollOption.PollOptionID == pollID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Poll Option ID already exists"})
			return
		}
	}

	if _, err := p.helper.JSONSet(cacheKey, ".", pollItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new PollOptions"})
		return
	}

	c.JSON(http.StatusOK, pollItem)
}

func (p *PollAPI) GetPoll(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Poll ID provided"})
		return
	}

	cacheKey := "poll:" + id
	pollBytes, err := p.helper.JSONGet(cacheKey, ".")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll in cache with id=" + cacheKey})
		return
	}

	var poll poll.Poll
	err = json.Unmarshal(pollBytes.([]byte), &poll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cached data seems to be wrong type"})
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (p *PollAPI) GetPollList(c *gin.Context) {
	var pollList []poll.Poll
	var pollItem poll.Poll

	pattern := "poll:*"
	ks, _ := p.client.Keys(p.context, pattern).Result()
	for _, key := range ks {
		err, _ := p.getItemFromRedis(key, &pollItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find poll in cache with id=" + key})
			return
		}
		pollList = append(pollList, pollItem)
	}

	c.JSON(http.StatusOK, pollList)
}

func (p *PollAPI) GetPollOptions(c *gin.Context) {
	id := c.Param("id")
	var thisPoll poll.Poll
	pattern := "poll:" + id
	err, _ := p.getItemFromRedis(pattern, &thisPoll)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll in cache with id=" + pattern})
		return
	}

	c.JSON(http.StatusOK, thisPoll.PollOptions)
}

func (p *PollAPI) GetPollOption(c *gin.Context) {
	id := c.Param("id")
	pollIDS := c.Param("pollid")
	if pollIDS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll option ID provided"})
		return
	}
	pollID64, err := strconv.ParseInt(pollIDS, 10, 32)
	if err != nil {
		log.Println("Error converting poll option id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	cacheKey := "poll:" + id

	var existingPoll poll.Poll
	err, existingPollOptions := p.getItemFromRedis(cacheKey, &existingPoll)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll in cache with id=" + cacheKey})
		return
	}

	pollID := uint(pollID64)

	for _, pollOption := range existingPollOptions {
		if pollOption.PollOptionID == pollID {
			c.JSON(http.StatusOK, pollOption)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll option in cache with id=" + pollIDS})
}

func (p *PollAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"status":             "ok",
			"version":            "3.0.0",
			"uptime":             300,
			"users_processed":    3000,
			"errors_encountered": 30,
		})
}

func (p *PollAPI) getItemFromRedis(key string, poll *poll.Poll) (error, []poll.PollOption) {
	itemObject, err := p.helper.JSONGet(key, ".")
	if err != nil {
		return err, poll.PollOptions
	}

	err = json.Unmarshal(itemObject.([]byte), poll)
	if err != nil {
		return err, poll.PollOptions
	}

	return nil, poll.PollOptions
}
