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

	vl, err := api.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r.GET("/voters", vl.GetVoterList)
	r.GET("/voters/:id", vl.GetVoter)
	r.POST("/voters/:id", vl.AddVoter)
	r.GET("/voters/:id/polls", vl.GetPolls)
	r.GET("/voters/:id/polls/:pollid", vl.GetPoll)
	r.POST("/voters/:id/polls/:pollid", vl.AddPoll)
	r.GET("/voters/health", vl.HealthCheck)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
