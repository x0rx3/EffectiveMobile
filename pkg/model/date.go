package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Date time.Time

// Для JSON
func (inst *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		*inst = Date(time.Time{})
		return nil
	}
	t, err := time.Parse("2006-01-02", s) // формат YYYY-MM-DD
	if err != nil {
		return err
	}
	*inst = Date(t)
	return nil
}

func (inst Date) MarshalJSON() ([]byte, error) {
	t := time.Time(inst)
	return []byte(fmt.Sprintf(`"%s"`, t.Format("2006-01-02"))), nil
}

// Для pgx
func (inst Date) Value() (driver.Value, error) {
	return time.Time(inst), nil // pgx умеет работать с time.Time
}

func (inst *Date) Scan(src interface{}) error {
	if src == nil {
		*inst = Date(time.Time{})
		return nil
	}
	switch t := src.(type) {
	case time.Time:
		*inst = Date(t)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Date", src)
	}
}
