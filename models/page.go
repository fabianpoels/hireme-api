package models

type Page struct {
	Index   int      `json:"index" bson:"index"`
	Key     string   `json:"key" bson:"key"`
	Guesses []string `json:"guesses" bson:"guesses"`
	Hints   int      `json:"attempts" bson:"attempts"`
}
