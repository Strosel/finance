package main

import (
	"context"
	"fmt"
	"math/rand"
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

func RandEventStub(n int) []Event {
	ev := []Event{}

	rand.Seed(time.Now().Unix())

	for i := 0; i < n; i++ {
		if rand.Float64() < 0.5 {
			ev = append(ev, Receipt{
				Datetime: time.Now().Add(time.Duration(-i) * time.Hour),
				Store:    fmt.Sprintf("Store%v", i),
				Products: []*Transaction{
					&Transaction{
						Name:     fmt.Sprintf("Prod%v", i),
						Category: fmt.Sprintf("Cat_%v", rand.Intn(5)),
						Sum:      rand.Intn(500),
					},
				},
			})
		} else {
			ev = append(ev, Transaction{
				Datetime: time.Now().Add(time.Duration(-i) * time.Hour),
				Name:     fmt.Sprintf("Name%v", i),
				Category: fmt.Sprintf("Cat_%v", rand.Intn(5)),
				Sum:      rand.Intn(500),
			})
		}
	}

	return ev
}

func GetEvents(start, end time.Time) []Event {
	trs := []Transaction{}
	res := []Receipt{}
	//? Timeout len
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

	//TODO SORT
	evs := []Event{}
	for _, t := range trs {
		evs = append(evs, t)
	}
	for _, r := range res {
		evs = append(evs, r)
	}
	return evs
}
