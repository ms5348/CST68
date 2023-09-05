#!/bin/bash
curl -H "Content-Type: application/json" -X GET http://localhost:1080/voters
curl -d '{ "voter_id": 1 }' -H "Content-Type: application/json" -X GET http://localhost:1080/voters/1 
curl -d '{ "voter_id": 3 }' -H "Content-Type: application/json" -X GET http://localhost:1080/voters/3/polls 
curl -d '{ "voter_id": 3, "vote_history": [{"poll_id": 1}] }' -H "Content-Type: application/json" -X GET http://localhost:1080/voters/3/polls/1
curl -H "Content-Type: application/json" -X GET http://localhost:1080/voters/health

curl -H "Content-Type: application/json" -X GET http://localhost:1082/poll
curl -d '{ "poll_id": 1 }' -H "Content-Type: application/json" -X GET http://localhost:1082/poll/1
curl -d '{ "poll_id": 2 }' -H "Content-Type: application/json" -X GET http://localhost:1082/poll/2/polloptions
#curl -d '{ "poll_id": 2, "poll_options": [{"option_id": 2}] }' -H "Content-Type: application/json" -X GET http://localhost:1082/poll/2/polloptions/2
curl -H "Content-Type: application/json" -X GET http://localhost:1082/poll/health

curl -H "Content-Type: application/json" -X GET http://localhost:1081/votes
curl -d '{ "vote_id": 4 }' -H "Content-Type: application/json" -X GET http://localhost:1081/votes/4 
curl -H "Content-Type: application/json" -X GET http://localhost:1081/votes/health
