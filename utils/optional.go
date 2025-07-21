package utils

import (
	"fmt"
)

var ErrOptionalHasNoValue = fmt.Errorf("Optional value does not contain a value.")

type Optional[T any] struct {
	valuePresent bool
	value        T
}

func SomeOptional[T any](val T) Optional[T] {
	return Optional[T]{
		value:        val,
		valuePresent: true,
	}
}

func NoneOptional[T any]() Optional[T] {
	return Optional[T]{
		valuePresent: false,
	}
}

func NewSomeOptional[T any](val T) *Optional[T] {
	o := new(Optional[T])

	o.value = val
	o.valuePresent = true

	return o
}

func NewNoneOptional[T any](val T) *Optional[T] {
	o := new(Optional[T])

	o.valuePresent = false

	return o
}

func (o Optional[T]) Value() (T, error) {
	if o.valuePresent {
		return o.value, nil
	} else {
		//log.Println(string(debug.Stack()))
		return o.value, ErrOptionalHasNoValue
	}
}

func OptionalMap[A any, B any](o Optional[A], f func(A) B) Optional[B] {
	if o.valuePresent {
		return SomeOptional(f(o.value))
	} else {
		return NoneOptional[B]()
	}
}
