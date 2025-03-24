package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type MilliTimestamp int64

func (ts MilliTimestamp) Value() (driver.Value, error) {
	return int64(ts), nil
}

func (ts *MilliTimestamp) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		*ts = MilliTimestamp(0)
		return nil
	case int64:
		*ts = MilliTimestamp(v)
		return nil
	default:
		return nil
	}
}

func (ts *MilliTimestamp) MarshalJSON() ([]byte, error) {
	if *ts == 0 {
		return json.Marshal(nil)
	}
	stamp := time.Unix(int64(*ts)/1000, int64(*ts)%1000*int64(time.Millisecond)).Format("2006-01-02 15:04:05.000")
	return json.Marshal(stamp)
}

func (ts *MilliTimestamp) UnmarshalJSON(b []byte) error {
	var s *string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	if s == nil {
		*ts = MilliTimestamp(0)
		return nil
	}
	tt, err := time.Parse("2006-01-02 15:04:05.000", *s)
	if err != nil {
		return err
	}
	*ts = MilliTimestamp(tt.UnixNano() / int64(time.Millisecond))
	return nil
}
