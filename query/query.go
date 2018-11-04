//Package query provides functionality to get information on elements in an array.
package query

import "reflect"

// ArrayContains returns true if the array passed in contains the value passed in.
func ArrayContains(val interface{}, array interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		slice := reflect.ValueOf(array)

		for i := 0; i < slice.Len(); i++ {
			if reflect.DeepEqual(val, slice.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}
