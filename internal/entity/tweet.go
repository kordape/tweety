// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

type TweetWithClassification struct {
	Text string `json:"body"`
	Fake int    `json:"fake"`
}

type Tweet struct {
	Text string
}
