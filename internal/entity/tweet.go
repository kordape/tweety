// Package entity defines main entities for business logic (services), database mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

type TweetWithClassification struct {
	Text      string `json:"body"`
	Fake      int    `json:"fake"`
	CreatedAt string `json:"created_at"`
}

type Tweet struct {
	CreatedAt string `json:"created_at"`
	Id        string `json:"id"`
	Text      string `json:"text"`
}
