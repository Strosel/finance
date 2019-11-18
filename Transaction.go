package main

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID       primitive.ObjectID `bson:"_id"`
	Datetime time.Time          `bson:"datetime,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Sum      int                `bson:"sum,omitempty"` // int of ören obvuiously
	Category string             `bson:"category,omitempty"`
	Note     string             `bson:"note,omitempty"`
}

func (t Transaction) GetTime() time.Time {
	return t.Datetime
}

func (t Transaction) GetName() string {
	return t.Name
}

func (t Transaction) GetSum() int {
	return t.Sum
}

func (t Transaction) GetCategory() string {
	return t.Category
}

func (t Transaction) GetType() string {
	id := fmt.Sprintf("T/" + t.ID.Hex())
	return idre.ReplaceAllString(id, "/0")
}