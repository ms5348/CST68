package api

import (
	"log"
	"net/http"
	"strconv"

	"final-assignment/votes"

	"github.com/gin-gonic/gin"
)

type VotesAPI struct {
	votes *votes.VoteCache
}

func NewVotesAPI() (*VotesAPI, error) {
	votesHandler, err := votes.New()
	if err != nil {
		return nil, err
	}

	return &VotesAPI{votes: votesHandler}, nil
}

func (v *VotesAPI) AddVotes(c *gin.Context) {
	var newVotes votes.Vote

	if err := c.ShouldBindJSON(&newVotes); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.votes.AddVote(newVotes); err != nil {
		log.Println("Error adding vote: ", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.JSON(http.StatusOK, newVotes)
}

func (v *VotesAPI) GetVotes(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	votes, err := v.votes.GetVote(uint(id64))
	if err != nil {
		log.Println("Vote not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, votes)
}

func (v *VotesAPI) GetVotesList(c *gin.Context) {
	voteList, err := v.votes.GetAllVotes()
	if err != nil {
		log.Println("Error Getting All Votes: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if voteList == nil {
		voteList = make([]votes.Vote, 0)
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
