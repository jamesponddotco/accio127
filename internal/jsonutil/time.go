package jsonutil

import (
	"encoding/json"
	"fmt"
	"time"
)

// Duration is a wrapper around time.Duration which supports JSON unmarshalling
// from a string.
type Duration time.Duration

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return fmt.Errorf("could not unmarshal duration as string: %w", err)
	}

	duration, err := time.ParseDuration(str)
	if err != nil {
		return fmt.Errorf("could not parse duration string: %w", err)
	}

	*d = Duration(duration)

	return nil
}
