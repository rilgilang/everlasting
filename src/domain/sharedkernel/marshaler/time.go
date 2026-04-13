package marshaler

import (
	"strconv"
	"time"
)

type JsonTime time.Time

type MyStruct struct {
	Date JsonTime
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t *JsonTime) UnmarshalJSON(s []byte) (err error) {
	q, err := strconv.ParseInt(string(s), 10, 64)

	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q, 0)
	return
}

func (t JsonTime) String() string { return time.Time(t).String() }
