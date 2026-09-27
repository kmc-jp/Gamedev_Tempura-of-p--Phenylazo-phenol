package Utils

import "reflect"

// reflect.TypeFor[T]() を使用せよ
func Typeof[T any]() reflect.Type {
	println("reflect.TypeFor[T]() を使用せよ")
	return reflect.TypeOf((*T)(nil)).Elem()
}
