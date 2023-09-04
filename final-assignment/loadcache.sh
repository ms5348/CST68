#!/bin/bash
curl -d '{ "voter_id": 1, "first_name": "John", "last_name": "Doe" }' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/1 
curl -d '{ "voter_id": 2, "first_name": "Jane", "last_name": "Doe"}' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/2 
curl -d '{ "voter_id": 3, "first_name": "Matt", "last_name": "S"}' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/3
curl -d '{ "voter_id": 3, "vote_history": [{"poll_id": 0}]}' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/3/polls/1	