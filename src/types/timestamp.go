package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Timestamp int

func (ts Timestamp) Value() (driver.Value, error) {
	return int64(ts), nil
}
func (ts *Timestamp) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		*ts = Timestamp(0)
		return nil
	case int64:
		*ts = Timestamp(v)
		return nil
	default:
		return nil
	}
}

func (ts *Timestamp) MarshalJSON() ([]byte, error) {
	if *ts == 0 {
		return json.Marshal(nil)
	}
	stamp := time.Unix(int64(*ts), 0).Format("2006-01-02 15:04:05")
	return json.Marshal(stamp)
}

func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	var s *string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	if s == nil {
		*ts = Timestamp(0)
		return nil
	}
	tt, err := time.Parse("2006-01-02 15:04:05", *s)
	if err != nil {
		return err
	}
	*ts = Timestamp(tt.Unix())
	return nil
}
