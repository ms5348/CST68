package main

import (
	"final-assignment/api"
	"flag"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

func main() {
	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	voterAPI, err := api.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pAPI, err := api.NewPollAPI()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	votesAPI, err := api.NewVotesAPI()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r.GET("/voters", voterAPI.GetVoterList)
	r.GET("/voters/:id", voterAPI.GetVoter)
	r.POST("/voters/:id", voterAPI.AddVoter)
	r.GET("/voters/:id/polls", voterAPI.GetPolls)
	r.GET("/voters/:id/polls/:pollid", voterAPI.GetPoll)
	r.POST("/voters/:id/polls/:pollid", voterAPI.AddPoll)
	r.GET("/voters/health", voterAPI.HealthCheck)

	r.GET("/poll", pAPI.GetPollList)
	r.GET("/poll/:id", pAPI.GetPoll)
	r.POST("/poll/:id", pAPI.AddPoll)
	r.GET("/poll/:id/polloptions", pAPI.GetPollOptions)
	//rPoll.GET("/poll/:id/polloptions/:pollid", pAPI.GetPollOption)
	//rPoll.POST("/poll/:id/polloptions/:pollid", pAPI.AddPollOption)
	r.GET("/poll/health", pAPI.HealthCheck)

	r.GET("/votes", votesAPI.GetVotesList)
	r.GET("/votes/:id", votesAPI.GetVotes)
	r.POST("/votes/:id", votesAPI.AddVotes)
	r.GET("/votes/health", votesAPI.HealthCheck)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
