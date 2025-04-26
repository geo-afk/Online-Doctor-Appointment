package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgText(value string) pgtype.Text {

	var text pgtype.Text
	text.String = value
	text.Valid = true
	return text
}

func ToPgInt(value int32) pgtype.Int4 {

	var number pgtype.Int4
	number.Int32 = value
	number.Valid = true
	return number
}
func ToPgTime(value time.Time) pgtype.Timestamp {

	var time pgtype.Timestamp
	time.Time = value
	time.Valid = true
	return time
}
