package helpers_datetime

import (
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"time"
)

func ParseStringGmt(timeString string) int64 {
	t, err := time.Parse(time.RFC3339, timeString)
	helpers_error.PanicIfError(err)
	return t.Unix()
}
func ParseTimestampGmt(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(time.DateTime)
}
