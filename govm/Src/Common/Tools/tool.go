package Tools

import (
	"fmt"
	"time"
)

func ToTimeStr(t time.Time) string {
	return fmt.Sprintf("%d:%d:%d", t.Hour(), t.Minute(), t.Second())
}

func ToDateStr(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%d-%d-%d", y, m, d)
}

func ToDateTimeStr(t time.Time) string {
	return fmt.Sprintf("%s %s", ToDateStr(t), ToTimeStr(t))
}
