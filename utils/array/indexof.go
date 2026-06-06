package array

import "reflect"

func IndexOf(arr interface{}, element interface{}) int {
	val := reflect.ValueOf(arr)
	if !val.IsValid() {
		return -1
	}
	k := val.Kind()
	if k != reflect.Slice && k != reflect.Array {
		return -1
	}
	for i := 0; i < val.Len(); i++ {
		if reflect.DeepEqual(val.Index(i).Interface(), element) {
			return i
		}
	}
	return -1
}
