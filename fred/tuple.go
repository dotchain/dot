// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import "reflect"

// ToTuple converts a slice of objects into an array
func ToTuple(objects []Object) interface{} {
	if len(objects) == 0 {
		return nil
	}
	t := reflect.ArrayOf(len(objects), reflect.TypeOf([]Object{}).Elem())
	result := reflect.New(t).Elem()
	for i, o := range objects {
		if o != nil {
			result.Index(i).Set(reflect.ValueOf(o))
		}
	}
	//reflect.Copy(result.Slice(0, len(objects)), reflect.ValueOf(objects))
	return result.Interface()
}

// FromTuple converts an array of objects into a slice
func FromTuple(v interface{}) []Object {
	if v == nil {
		return nil
	}
	r := reflect.ValueOf(v)
	result := make([]Object, r.Len())
	for kk := 0; kk < r.Len(); kk++ {
		if v := r.Index(kk).Interface(); v != nil {
			result[kk] = v.(Object)
		}
	}
	return result
}
