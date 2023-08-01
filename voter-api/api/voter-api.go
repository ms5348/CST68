package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"voter-api-starter/voter"

	"github.com/gin-gonic/gin"
)

type VoterApi struct {
	voterList voter.VoterList
}

func NewVoterApi() *VoterApi {
	return &VoterApi{
		voterList: voter.VoterList{
			Voters: make(map[uint]voter.Voter),
		},
	}
}

func (v *VoterApi) AddVoter(c *gin.Context) {
	var newVoter voter.Voter

	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Note: if the body of the POST request contains a VoterID key value, the following stores newVoter
	//		with the Post parameter (:id). The actual POST would use the key value though.
	idS := c.Param("id")
	id64, _ := strconv.ParseInt(idS, 10, 32)

	v.voterList.Voters[uint(id64)] = *voter.NewVoter(uint(id64), newVoter.FirstName, newVoter.LastName)

	c.JSON(http.StatusOK, newVoter)
}

func (v *VoterApi) AddPoll(c *gin.Context) {
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
	voter := v.voterList.Voters[uint(id64)]
	voter.AddPoll(uint(pollID64))
	v.voterList.Voters[uint(id64)] = voter
}

func (v *VoterApi) GetVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter := v.voterList.Voters[uint(id64)]

	c.JSON(http.StatusOK, voter)
}

func (v *VoterApi) GetVoterJson(voterID uint) string {
	voter := v.voterList.Voters[voterID]
	return voter.ToJson()
}

func (v *VoterApi) GetVoterList(c *gin.Context) {
	c.JSON(http.StatusOK, v.voterList)
}

func (v *VoterApi) GetVoterListJson() string {
	b, _ := json.Marshal(v.voterList)
	return string(b)
}

func (v *VoterApi) GetPolls(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	polls := v.voterList.Voters[uint(id64)].VoteHistory

	c.JSON(http.StatusOK, polls)
}

func (v *VoterApi) GetPoll(c *gin.Context) {
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

	poll := v.voterList.Voters[uint(id64)].VoteHistory[uint(pollID64)]

	c.JSON(http.StatusOK, poll)
}

func (v *VoterApi) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"status":             "ok",
			"version":            "1.0.0",
			"uptime":             100,
			"users_processed":    1000,
			"errors_encountered": 10,
		})
}
