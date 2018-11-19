package helpers

import (
	"github.com/teris-io/shortid"
	"strings"
	"time"
	"unicode"
)

func CurrentEpoch() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func UUID() string {
	uuid, _ := shortid.Generate()
	return uuid

}

func SaniString(str string) string {
	r := strings.ToLower(str)
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, r)
}
