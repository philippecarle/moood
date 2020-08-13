package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Entry defines the structure of an Entry when submitted by a user
type Entry struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	Content   string             `json:"content" binding:"required" bson:"content"`
	UserID    primitive.ObjectID `json:"-" bson:"userId"`
}

// SentimentallyAnalyzedEntry defines the structure of an Entry when processed by the sentiment analysis worker
type SentimentallyAnalyzedEntry struct {
	Entry     `bson:",inline"`
	Sentences []Sentence `json:"sentences" bson:"sentences"`
	Score     Score      `json:"score" bson:"score"`
}

// Sentence describes a sentence with the sentiment analysis results
type Sentence struct {
	Sentence string `json:"sentence" bson:"sentence"`
	Score    Score  `json:"score" bson:"score"`
}

// Score gives the metrics returned by the sentiment analysis worker
type Score struct {
	Negative float32 `json:"neg" bson:"negative"`
	Neutral  float32 `json:"neu" bson:"neutral"`
	Positive float32 `json:"pos" bson:"positive"`
	Compound float32 `json:"compound" bson:"compound"`
}

// FullEntry defines the structure of a full entry
type FullEntry struct {
	SentimentallyAnalyzedEntry `bson:",inline"`
}
