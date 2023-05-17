package model

// Counter represents an access counter.
type Counter struct {
	Count uint64 `json:"count"`
}

// NewCounter creates a new Counter.
func NewCounter(count uint64) *Counter {
	if count < 1 {
		count = 1
	}

	return &Counter{
		Count: count,
	}
}
