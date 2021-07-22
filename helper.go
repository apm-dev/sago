package sago

import (
	"reflect"
	"strings"
)

func structName(s interface{}) string {
	str, ok := s.(string)
	if ok {
		return str
	}
	n := strings.Split(reflect.TypeOf(s).String(), ".")
	return n[len(n)-1]
}
