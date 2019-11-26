package finance

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

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

func GetBudget(db *mongo.Collection, to time.Duration, t time.Time) (Budget, error) {
	bgt := Budget{}
	ctx, _ := context.WithTimeout(context.Background(), to)
	res := db.FindOne(ctx, bson.M{
		//time is modified to accomodate fo wierdness in mongodb
		"start": bson.M{
			"$lte": t.AddDate(0, 0, 1),
		},
		"end": bson.M{
			"$gte": t.AddDate(0, 0, -1),
		},
	})
	err := res.Decode(&bgt)

	return bgt, err
}
