#!/bin/bash
curl -d '{ "voter_id": 1, "first_name": "John", "last_name": "Doe" }' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/1 
curl -d '{ "voter_id": 2, "first_name": "Jane", "last_name": "Doe"}' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/2 
curl -d '{ "voter_id": 3, "first_name": "Matt", "last_name": "S"}' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/3
curl -d '{ "voter_id": 3, "vote_history": [{"poll_id": 1}]}' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/3/polls/1	

curl -d '{ "poll_id": 1, "poll_title": "Favorite Pet", "poll_question": "What type of pet do you like best?" }' -H "Content-Type: application/json" -X POST http://localhost:1082/poll/1 
curl -d '{ "poll_id": 2, "poll_title": "Favorite Color", "poll_question": "What color do you like best?" }' -H "Content-Type: application/json" -X POST http://localhost:1082/poll/2
curl -d '{ "poll_id": 3, "poll_title": "Favorite Car", "poll_question": "What type of car do you like best?" }' -H "Content-Type: application/json" -X POST http://localhost:1082/poll/3 
#curl -d '{ "poll_id": 1, "poll_options": [{"option_id": 1, "option_text": "Dog"}]}' -H "Content-Type: application/json" -X POST http://localhost:1082/poll/1/polloptions/1	
#curl -d '{ "poll_id": 1, "poll_options": [{"option_id": 2, "option_text": "Cat"}]}' -H "Content-Type: application/json" -X POST http://localhost:1082/poll/1/polloptions/2	
#curl -d '{ "poll_id": 2, "poll_options": [{"option_id": 1, "option_text": "Red"}]}' -H "Content-Type: application/json" -X POST http://localhost:1082/poll/2/polloptions/1	
#curl -d '{ "poll_id": 2, "poll_options": [{"option_id": 2, "option_text": "Blue"}]}' -H "Content-Type: application/json" -X POST http://localhost:1082/poll/2/polloptions/2	

curl -d '{ "vote_id": 1, "voter_id": 1, "poll_id": 1, "vote_value": 1 }' -H "Content-Type: application/json" -X POST http://localhost:1081/votes/1 
curl -d '{ "vote_id": 2, "voter_id": 1, "poll_id": 2, "vote_value": 2 }' -H "Content-Type: application/json" -X POST http://localhost:1081/votes/2 
curl -d '{ "vote_id": 3, "voter_id": 2, "poll_id": 1, "vote_value": 2 }' -H "Content-Type: application/json" -X POST http://localhost:1081/votes/3 
curl -d '{ "vote_id": 4, "voter_id": 2, "poll_id": 2, "vote_value": 1 }' -H "Content-Type: application/json" -X POST http://localhost:1081/votes/4 
