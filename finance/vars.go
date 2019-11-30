package finance

import (
	"regexp"
	"strings"
)

var (
	zre     = regexp.MustCompile(`/0+`)
	savestr = regexp.MustCompile("(spar|spara|sparande|save|saving|savings)")
)

func IsSavings(obj interface{}) bool {
	switch o := obj.(type) {
	case Transaction:
		return savestr.MatchString(strings.ToLower(o.Name))
	case string:
		return savestr.MatchString(strings.ToLower(o))
	default:
		return false
	}
}
