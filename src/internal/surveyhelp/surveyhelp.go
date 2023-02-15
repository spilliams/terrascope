// Package surveyhelp provides helper functions for use with the package
// github.com/AlecAivazis/survey/v2.
package surveyhelp

import (
	"reflect"
	"strings"
)

// SplitTransformer returns a transform function that takes the survey answer
// and splits it on a given character.
func SplitTransformer(char string) func(ans interface{}) interface{} {
	return func(ans interface{}) interface{} {
		if isZero(reflect.ValueOf(ans)) {
			return []string{}
		}

		s, ok := ans.(string)
		if !ok {
			return []string{}
		}

		return strings.Split(s, char)
	}
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
