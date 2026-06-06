package array

import "reflect"

func Contain(arr interface{}, element interface{}) bool {
	v := reflect.ValueOf(arr)
	k := v.Kind()
	if k != reflect.Slice && k != reflect.Array {
		return false
	}
	for i := 0; i < v.Len(); i++ {
		if reflect.DeepEqual(v.Index(i).Interface(), element) {
			return true
		}
	}
	return false
}
