package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Event interface {
	GetTime() time.Time
	GetName() string
	GetSum() int
	GetType() string
}

func RandEventStub(n int) []Event {
	ev := []Event{}

	rand.Seed(time.Now().Unix())

	for i := 0; i < n; i++ {
		if rand.Float64() < 0.5 {
			ev = append(ev, Receipt{
				Store: fmt.Sprintf("Store%v", i),
			})
		} else {
			ev = append(ev, Transaction{
				Datetime: time.Now().Add(time.Duration(-i) * time.Hour),
				Name:     fmt.Sprintf("Name%v", i),
				Sum:      rand.Intn(500),
			})
		}
	}

	return ev
}
