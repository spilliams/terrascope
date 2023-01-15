package surveyhelp

import (
	"reflect"
	"strings"
)

func SplitTransformer(ans interface{}) interface{} {
	if isZero(reflect.ValueOf(ans)) {
		return []string{}
	}

	s, ok := ans.(string)
	if !ok {
		return []string{}
	}

	return strings.Split(s, " ")
}

// isZero returns true if the passed value is the zero object
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	}

	// compare the types directly with more general coverage
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
