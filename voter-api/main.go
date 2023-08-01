package main

import (
	"flag"
	"fmt"
	"voter-api-starter/api"
	"voter-api-starter/voter"

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

	v := voter.NewVoter(1, "John", "Doe")
	v.AddPoll(1)
	//v.AddPoll(2)
	//v.AddPoll(3)
	//v.AddPoll(4)

	//b, _ := json.Marshal(v)
	//fmt.Println(string(b))
	vl := api.NewVoterApi()
	//vl.AddVoter(1, "John", "Doe")
	//vl.AddPoll(1, 1)
	//vl.AddPoll(1, 2)
	//vl.AddVoter(2, "Jane", "Doe")
	//vl.AddPoll(2, 1)
	//vl.AddPoll(2, 2)

	fmt.Println("------------------------")
	fmt.Println(vl.GetVoterJson(1))
	fmt.Println("------------------------")
	fmt.Println(vl.GetVoterJson(2))
	fmt.Println("------------------------")
	fmt.Println(vl.GetVoterListJson())
	fmt.Println("------------------------")

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
