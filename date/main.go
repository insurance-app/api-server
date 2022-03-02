package date

import (
	"encoding/json"
	"strings"
	"time"
)

type Date time.Time

func (j *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = Date(t)
	return nil
}

func (j Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.String())
}

func (j Date) String() string {
	return time.Time(j).Format("2006-01-02")
}
