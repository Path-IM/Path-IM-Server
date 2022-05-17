package global

import (
	"database/sql"
	"database/sql/driver"
	timeUtils "github.com/showurl/Zero-IM-Server/common/utils/time"
	"time"
)

type ModelTime time.Time
type ModelNullTime sql.NullTime

type ModelDate time.Time

func (myT ModelDate) MarshalText() ([]byte, error) {
	s := time.Time(myT).Format("2006-01-02")
	return []byte(s), nil
}

func (myT *ModelDate) UnmarshalText(text []byte) (err error) {
	t := (*time.Time)(myT)
	*t, err = time.ParseInLocation("2006-01-02", string(text), time.Local)
	return
}

func (myT ModelDate) String() string {
	return time.Time(myT).Format(timeUtils.DateFormat)
}

func (myT ModelDate) Value() (driver.Value, error) {
	return time.Time(myT), nil
}

func (myT ModelDate) Time() time.Time {
	return time.Time(myT)
}

func (myT ModelTime) MarshalText() ([]byte, error) {
	s := time.Time(myT).Format("2006-01-02 15:04:05.000")
	return []byte(s), nil
}

func (myT *ModelTime) UnmarshalText(text []byte) (err error) {
	t := (*time.Time)(myT)
	*t, err = time.ParseInLocation("2006-01-02 15:04:05.000", string(text), time.Local)
	return
}

func (myT ModelTime) String() string {
	return time.Time(myT).Format(timeUtils.TimeFormat)
}

func (myT ModelTime) PageIndex() int64 {
	return myT.Time().UnixNano() / 1e6
}

func (myT ModelTime) Value() (driver.Value, error) {
	return time.Time(myT), nil
}

func (myT ModelTime) Time() time.Time {
	return time.Time(myT)
}

func (myT ModelNullTime) MarshalText() ([]byte, error) {
	s := ""
	nullTime := sql.NullTime(myT)
	if nullTime.Valid {
		s = nullTime.Time.Format("2006-01-02 15:04:05.000")
	}
	return []byte(s), nil
}

func (myT *ModelNullTime) UnmarshalText(text []byte) error {
	t := (*sql.NullTime)(myT)
	if string(text) != "" {
		tm, err := time.ParseInLocation("2006-01-02 15:04:05.000", string(text), time.Local)
		if err != nil {
			return err
		}
		*t = sql.NullTime{
			Valid: true,
			Time:  tm,
		}
	}
	return nil
}

func (myT ModelNullTime) String() string {
	nullTime := sql.NullTime(myT)
	if nullTime.Valid {
		return nullTime.Time.Format(timeUtils.TimeFormat)
	} else {
		return ""
	}
}

func (myT ModelNullTime) Value() (driver.Value, error) {
	nullTime := sql.NullTime(myT)
	if nullTime.Valid {
		return nullTime.Time, nil
	}
	return nil, nil
}

func (myT *ModelNullTime) Scan(value interface{}) error {
	if value == nil {
		myT.Time, myT.Valid = time.Time{}, false
		return nil
	}
	myT.Valid = true
	s := value.(time.Time)
	d := &myT.Time
	*d = s
	return nil
}

func (myT ModelNullTime) NullTime() sql.NullTime {
	return sql.NullTime(myT)
}
