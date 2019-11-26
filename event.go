package main

import (
	"context"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Event interface {
	GetTime() time.Time
	GetName() string
	GetSum() int
	GetSumS() string
	GetCategory() string
	GetType() string
}

func GetEvents(start, end time.Time) []Event {
	trs := []Transaction{}
	res := []Receipt{}
	ctx, _ := context.WithTimeout(context.Background(), dTimeout)
	curs, err := db.Collection(tDb).Find(ctx, &bson.D{bson.E{
		Key: "datetime",
		Value: bson.D{
			bson.E{Key: "$gte", Value: start},
			bson.E{Key: "$lte", Value: end},
		}}})
	if err != nil {
		ui.SetWidget(NewErrorView(err))
		ui.SetFocusChain(nil)
	}
	defer curs.Close(ctx)
	err = curs.All(ctx, &trs)
	if err != nil {
		ui.SetWidget(NewErrorView(err))
		ui.SetFocusChain(nil)
	}

	curs, err = db.Collection(rDb).Find(ctx, &bson.D{bson.E{
		Key: "datetime",
		Value: bson.D{
			bson.E{Key: "$gte", Value: start},
			bson.E{Key: "$lte", Value: end},
		}}})
	if err != nil {
		ui.SetWidget(NewErrorView(err))
		ui.SetFocusChain(nil)
	}
	defer curs.Close(ctx)
	err = curs.All(ctx, &res)
	if err != nil {
		ui.SetWidget(NewErrorView(err))
		ui.SetFocusChain(nil)
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

	return evs
}
