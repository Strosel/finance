package main

import "time"

type Receipt struct {
	Datetime time.Time
	Store    string
	Products []Transaction
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

func (r Receipt) GetType() string {
	return "Receipt"
}
