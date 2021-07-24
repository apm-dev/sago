package common

import (
	"reflect"
	"strings"
)

func StructName(s interface{}) string {
	str, ok := s.(string)
	if ok {
		return str
	}
	n := strings.Split(reflect.TypeOf(s).String(), ".")
	return n[len(n)-1]
}
