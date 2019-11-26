package main

import (
	"context"
	"time"

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

func GetBudget(t time.Time) Budget {
	bgt := Budget{}
	ctx, _ := context.WithTimeout(context.Background(), dTimeout)
	res := db.Collection(bDb).FindOne(ctx, bson.M{
		//time is modified to accomodate fo wierdness in mongodb
		"start": bson.M{
			"$lte": t.AddDate(0, 0, 1),
		},
		"end": bson.M{
			"$gte": t.AddDate(0, 0, -1),
		},
	})
	res.Decode(&bgt)

	return bgt
}
