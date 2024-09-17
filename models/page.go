package models

type Page struct {
	Guesses []string `json:"guesses" bson:"guesses"`
	Hints   int      `json:"hints" bson:"hints"`
}
