package models

type Score struct {
	SessionId string `json:"sessionId" bson:"sessionId"`
	Username  string `json:"username" bson:"username"`
	Score     int    `json:"score" bson:"score"`
	Guesses   int    `json:"guesses" bson:"guesses"`
	Hints     int    `json:"hints" bson:"hints"`
}
