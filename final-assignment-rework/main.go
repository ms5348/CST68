package main

import (
	"final-assignment/api"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag    string
	portFlag    uint
	cacheURL    string
	voterAPIURL string
	pollAPIURL  string
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.StringVar(&voterAPIURL, "voterapi", "http://localhost:1080", "Default endpoint for voter API")
	flag.StringVar(&pollAPIURL, "pollapi", "http://localhost:1082", "Default endpoint for poll API")
	flag.StringVar(&cacheURL, "c", "0.0.0.0:6379", "Default cache location")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

func envVarOrDefault(envVar string, defaultVal string) string {
	envVal := os.Getenv(envVar)
	if envVal != "" {
		return envVal
	}
	return defaultVal
}

func setupParms() {
	processCmdLineFlags()

	cacheURL = envVarOrDefault("REDIS_URL", cacheURL)
	voterAPIURL = envVarOrDefault("VOTER_API_URL", voterAPIURL)
	pollAPIURL = envVarOrDefault("POLL_API_URL", pollAPIURL)
	hostFlag = envVarOrDefault("RLAPI_HOST", hostFlag)

	pfNew, err := strconv.Atoi(envVarOrDefault("VOTES_API_PORT", fmt.Sprintf("%d", portFlag)))
	if err == nil {
		portFlag = uint(pfNew)
	}

}

func main() {
	setupParms()
	log.Println("Init/cacheURL: " + cacheURL)
	log.Println("Init/voterAPIURL: " + voterAPIURL)
	log.Println("Init/pollAPIURL: " + pollAPIURL)
	log.Println("Init/hostFlag: " + hostFlag)
	log.Printf("Init/portFlag: %d", portFlag)

	r := gin.Default()
	r.Use(cors.Default())

	voterAPI, err := api.NewVoterAPI(cacheURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pAPI, err := api.NewPollAPI(cacheURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	votesAPI, err := api.NewVotesAPI(cacheURL, voterAPIURL, pollAPIURL)
	if err != nil {
		panic(err)
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
	r.GET("/poll/:id/polloptions/:pollid", pAPI.GetPollOption)
	r.POST("/poll/:id/polloptions/:pollid", pAPI.AddPollOption)
	r.GET("/poll/health", pAPI.HealthCheck)

	r.GET("/votes", votesAPI.GetVotesList)
	r.GET("/votes/:id", votesAPI.GetVotes)
	r.POST("/votes/:id", votesAPI.AddVote)
	r.GET("/votes/health", votesAPI.HealthCheck)
	r.GET("/votes/:id/:idx", votesAPI.GetItemFromVote)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
