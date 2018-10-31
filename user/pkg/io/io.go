package io

import (
	"encoding/json"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	Cedula       string        `json:"cedula" bson:"cedula"`
	Name         string        `json:"name"  bson:"name"`
	LunchroomID  string        `json:"lunchroom_id"  bson:"lunchroom_id"`
	ActiveTicket string        `json:"active_ticket"  bson:"active_ticket"`
}

func (t User) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		return "unsupported value type"
	}
	return string(b)
}
