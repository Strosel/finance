package finance

import "time"

type Income struct {
	Sum  int       `bson:"sum,omitempty"`
	Date time.Time `bson:"date,omitempty"`
}

func (i *Income) Received() bool {
	return time.Now().Before(i.Date)
}
