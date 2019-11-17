package main

import "time"

type Transaction struct {
	Datetime time.Time
	Name     string
	Sum      int // int of ören obvuiously
	Category string
	Note     string
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
	return "T"
}
