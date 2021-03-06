package lrdd

import "reflect"

var rowType = reflect.TypeOf((*Row)(nil))

func From(values interface{}) (rows []*Row) {
	inputVal := reflect.ValueOf(values)
	switch inputVal.Kind() {
	case reflect.Slice, reflect.Array:
		if inputVal.Type().Elem() == rowType {
			return values.([]*Row)
		}
		for i := 0; i < inputVal.Len(); i++ {
			rows = append(rows, Value(inputVal.Index(i).Interface()))
		}
	case reflect.Map:
		iter := inputVal.MapRange()
		for iter.Next() {
			k := iter.Key().String()
			v := iter.Value()
			if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
				for i := 0; i < v.Len(); i++ {
					rows = append(rows, KeyValue(k, v.Index(i).Interface()))
				}
			} else {
				rows = append(rows, KeyValue(k, v.Interface()))
			}
		}
	default:
		rows = append(rows, Value(values))
	}
	return
}
