package finance

import (
	"context"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Event interface {
	GetTime() time.Time
	GetName() string
	GetSum() int
	GetSumS() string
	GetCategory() string
	GetType() string
}

func GetEvents(db *mongo.Collection, to time.Duration, start, end time.Time) ([]Event, error) {
	trs := []Transaction{}
	res := []Receipt{}
	ctx, _ := context.WithTimeout(context.Background(), to)
	curs, err := db.Find(ctx, &bson.D{bson.E{
		Key: "datetime",
		Value: bson.D{
			bson.E{Key: "$gte", Value: start},
			bson.E{Key: "$lte", Value: end},
		}}})
	if err != nil {
		return nil, err
	}
	defer curs.Close(ctx)
	err = curs.All(ctx, &trs)
	if err != nil {
		return nil, err
	}

	curs, err = db.Find(ctx, &bson.D{bson.E{
		Key: "datetime",
		Value: bson.D{
			bson.E{Key: "$gte", Value: start},
			bson.E{Key: "$lte", Value: end},
		}}})
	if err != nil {
		return nil, err
	}
	defer curs.Close(ctx)
	err = curs.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	evs := []Event{}
	for _, t := range trs {
		evs = append(evs, t)
	}
	for _, r := range res {
		evs = append(evs, r)
	}

	sort.Slice(evs, func(i, j int) bool {
		return evs[i].GetTime().After(evs[j].GetTime())
	})

	return evs, err
}
