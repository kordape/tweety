// Package entity defines main entities for business logic (services), database mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import "time"

type TweetWithClassification struct {
	Text string `json:"body"`
	Fake int    `json:"fake"`
}

type Tweet struct {
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
	Text      string    `json:"text"`
}
