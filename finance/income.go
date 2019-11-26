package finance

import "time"

type Income struct {
	Sum      int       `bson:"sum,omitempty"`
	Received bool      `bson:"received,omitempty"`
	Date     time.Time `bson:"date,omitempty"`
}
