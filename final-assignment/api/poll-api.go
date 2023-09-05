package api

import (
	"log"
	"net/http"
	"strconv"

	"final-assignment/poll"

	"github.com/gin-gonic/gin"
)

type PollAPI struct {
	poll *poll.PollCache
}

func NewPollAPI() (*PollAPI, error) {
	pollHandler, err := poll.New()
	if err != nil {
		return nil, err
	}

	return &PollAPI{poll: pollHandler}, nil
}

func (p *PollAPI) AddPoll(c *gin.Context) {
	var newPoll poll.Poll

	if err := c.ShouldBindJSON(&newPoll); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := p.poll.AddPoll(newPoll); err != nil {
		log.Println("Error adding poll: ", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.JSON(http.StatusOK, newPoll)
}

/*
func (p *PollAPI) AddPoll(c *gin.Context) {
	var voter voter.Voter

	if err := c.ShouldBindJSON(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := p.voter.AddPoll(voter); err != nil {
		log.Println("Error adding poll: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, voter)
}*/

func (p *PollAPI) GetPoll(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	poll, err := p.poll.GetPoll(uint(id64))
	if err != nil {
		log.Println("Poll not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (p *PollAPI) GetPollList(c *gin.Context) {
	pollList, err := p.poll.GetAllPolls()
	if err != nil {
		log.Println("Error Getting All Polls: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if pollList == nil {
		pollList = make([]poll.Poll, 0)
	}

	c.JSON(http.StatusOK, pollList)
}

func (p *PollAPI) GetPollOptions(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollOptions, err := p.poll.GetPollOptions(uint(id64))
	if err != nil {
		log.Println("Poll not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, pollOptions)
}

/*
func (v *VoterAPI) GetPoll(c *gin.Context) {
	idS := c.Param("id")
	pollIDS := c.Param("pollid")
	id64, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	pollID64, err := strconv.ParseInt(pollIDS, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	poll, err := v.voter.GetPoll(uint(id64), uint(pollID64))
	if err != nil {
		log.Println("Voter or Poll not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, poll)
}*/

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
