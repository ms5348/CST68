#!/bin/bash
curl -d '{ "voter_id": 10 }' -H "Content-Type: application/json" -X GET http://localhost:1080/voters/10 
curl -d '{ "voter_id": 1 }' -H "Content-Type: application/json" -X GET http://localhost:1080/voters/1/polls #voter 1 does not have any polls stored but this and polloptions works right now
curl -d '{ "voter_id": 3, "vote_history": [{"poll_id": 10}] }' -H "Content-Type: application/json" -X GET http://localhost:1080/voters/3/polls/10

curl -d '{ "poll_id": 10 }' -H "Content-Type: application/json" -X GET http://localhost:1082/poll/10
curl -d '{ "poll_id": 3 }' -H "Content-Type: application/json" -X GET http://localhost:1082/poll/3/polloptions
#curl -d '{ "poll_id": 2, "poll_options": [{"option_id": 20}] }' -H "Content-Type: application/json" -X GET http://localhost:1082/poll/2/polloptions/20

curl -d '{ "vote_id": 40 }' -H "Content-Type: application/json" -X GET http://localhost:1081/votes/40 
