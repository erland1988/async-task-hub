package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Customtime time.Time

func (ts Customtime) Value() (driver.Value, error) {
	return time.Time(ts), nil
}

func (ts *Customtime) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if ok {
		*ts = Customtime(t)
		return nil
	}
	return nil
}

func (ts *Customtime) MarshalJSON() ([]byte, error) {
	stamp := time.Time(*ts).Format("2006-01-02 15:04:05")
	return json.Marshal(stamp)
}

func (ts *Customtime) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*ts = Customtime(t)
	return nil
}
