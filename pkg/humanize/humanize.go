package humanize

import (
	"fmt"
	"reflect"

	"github.com/dustin/go-humanize"
)

// Interface transform interface to string
func Interface(value interface{}) string {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		rValue := reflect.ValueOf(value)
		if rValue.IsNil() {
			return ""
		}

		value = rValue.Elem()
	}

	if reflect.TypeOf(value).Kind() == reflect.Map {
		return fmt.Sprintf("%d", reflect.ValueOf(value).Len())
	}
	if reflect.TypeOf(value).Kind() == reflect.Array || reflect.TypeOf(value).Kind() == reflect.Slice {
		return fmt.Sprintf("%d", reflect.ValueOf(value).Len())
	}

	switch value := value.(type) {
	case float64:
		return humanize.Ftoa(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}
