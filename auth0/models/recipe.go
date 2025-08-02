package models

import "time"

type Recipe struct {
	ID interface{} `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Tags []string `json:"tags" bson:"tags"`
	Ingredients []string `json:"ingredients" bson:"ingredients"`
	Instructions []string `json:"instructions" bson:"instructions"`
	PublishedAt time.Time `json:"publishedAt" bson:"publishedAt"`
}