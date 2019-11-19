package main

import (
	"context"
	"time"

	"github.com/strosel/noerr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Budget struct {
	ID       primitive.ObjectID `bson:"_id"`
	Start    time.Time          `bson:"start,omitempty"`
	End      time.Time          `bson:"end,omitempty"`
	Spending map[string]int     `bson:"spending,omitempty"`
	Income   map[string]Income  `bson:"income,omitempty"`
}

func GetBudget() Budget {
	bgt := Budget{}
	//? Timeout len
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	res := db.Collection("testB").FindOne(ctx, &bson.D{})
	noerr.Panic(res.Decode(&bgt))

	//TODO SORT
	return bgt
}
