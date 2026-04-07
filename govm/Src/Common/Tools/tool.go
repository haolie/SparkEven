package Tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"SparkEven/govm/Src/Common/Log"
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

func StrToCodeNum(str string) int {
	str = strings.Replace(str, ".SH", "", -1)
	str = strings.Replace(str, ".SZ", "", -1)
	str = strings.Replace(str, ".BJ", "", -1)
	temp, _ := strconv.Atoi(str)
	if temp < 1000000 {
		temp += 1000000
	}
	return temp
}

func GetCodeType(code int) string {
	code = code / 100000
	if code == 16 {
		return "sh"
	}

	if code == 10 || code == 13 {
		return "sz"
	}

	return "bj"
}

func ConvertCodeStr(code int) string {
	per := GetCodeType(code)
	temp := strconv.Itoa(code)[1:]
	return per + temp
}

func GetSecondsFromStr(timestr string) (int, bool) {
	h, m, s, success := ParseTimeStr(timestr)
	if success {
		return s + m*60 + h*60*60, true
	}

	return -1, false
}

func ParseTimeStr(timestr string) (h, m, s int, success bool) {
	strs := strings.Split(timestr, ":")
	if len(strs) == 3 {
		s = strToInt(strs[2])
		m = strToInt(strs[1])
		h = strToInt(strs[0])

		success = true
	}

	return h, m, s, success
}

func strToInt(str string) int {
	if len(str) == 1 {
		h, _ := strconv.Atoi(str)
		return h
	}

	t, _ := strconv.Atoi(str[0:1])
	m, _ := strconv.Atoi(str[1:])
	return t*10 + m
}

func GetTimeDate(t time.Time) time.Time {
	add := time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second
	return t.Add(-add)
}

func Max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func LoopCtx(ctx context.Context, fn func() bool, doneCb func()) bool {
	for {
		select {
		case <-ctx.Done():
			if doneCb != nil {
				doneCb()
			}

			return true

		default:
			if !fn() {
				return false
			}
		}
	}
}

func Wait(ctx context.Context, sec int, node string) {
	Log.Warn(fmt.Sprintf("%s wait %d", node, sec))
	w, _ := context.WithTimeout(ctx, time.Second*time.Duration(sec))
	<-w.Done()
}

func HunUpAfter(ctx context.Context, sec int, node string, fn func()) {
	w, _ := context.WithTimeout(ctx, time.Second*time.Duration(sec))
	fn()
	//Log.Warn(fmt.Sprintf("%s HunUpAfter %d", node, sec))
	<-w.Done()
}
