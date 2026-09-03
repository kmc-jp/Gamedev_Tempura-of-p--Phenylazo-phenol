package Utils

import "reflect"

func Typeof[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
