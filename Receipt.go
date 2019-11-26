package main

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Receipt struct {
	ID       primitive.ObjectID `bson:"_id"`
	Datetime time.Time          `bson:"datetime,omitempty"`
	Store    string             `bson:"store,omitempty"`
	Products []Transaction      `bson:"products,omitempty"`
}

func (r Receipt) GetTime() time.Time {
	return r.Datetime
}

func (r Receipt) GetName() string {
	return r.Store
}

func (r Receipt) GetSum() int {
	sum := 0
	for _, t := range r.Products {
		sum += t.Sum
	}
	return sum
}

func (r Receipt) GetSumS() string {
	return fmt.Sprintf("%8.2f", float64(r.GetSum())/100.)
}

func (r Receipt) GetCategory() string {
	return ""
}

func (r Receipt) GetType() string {
	id := fmt.Sprintf("R/" + r.ID.Hex())
	return idre.ReplaceAllString(id, "/0")
}
