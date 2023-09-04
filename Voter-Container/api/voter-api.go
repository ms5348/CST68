package api

import (
	"log"
	"net/http"
	"strconv"

	"Voter-Container/voter"

	"github.com/gin-gonic/gin"
)

type VoterAPI struct {
	voter *voter.VoterCache
}

func New() (*VoterAPI, error) {
	voterHandler, err := voter.New()
	if err != nil {
		return nil, err
	}

	return &VoterAPI{voter: voterHandler}, nil
}

func (v *VoterAPI) AddVoter(c *gin.Context) {
	var newVoter voter.Voter

	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.voter.AddVoter(newVoter); err != nil {
		log.Println("Error adding voter: ", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.JSON(http.StatusOK, newVoter)
}

func (v *VoterAPI) AddPoll(c *gin.Context) {
	var voter voter.Voter

	if err := c.ShouldBindJSON(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.voter.AddPoll(voter); err != nil {
		log.Println("Error adding poll: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, voter)
}

func (v *VoterAPI) GetVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, err := v.voter.GetVoter(uint(id64))
	if err != nil {
		log.Println("Voter not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter)
}

func (v *VoterAPI) GetVoterList(c *gin.Context) {
	voterList, err := v.voter.GetAllVoters()
	if err != nil {
		log.Println("Error Getting All Voters: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if voterList == nil {
		voterList = make([]voter.Voter, 0)
	}

	c.JSON(http.StatusOK, voterList)
}

func (v *VoterAPI) GetPolls(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	polls, err := v.voter.GetPolls(uint(id64))
	if err != nil {
		log.Println("Voter not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, polls)
}

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
