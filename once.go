package sync

import (
	"fmt"
	"reflect"
	"sync"
)

// OnceFunc returns a function that invokes fn only once.
// The returned function may be called concurrently.
func OnceFunc(fn func()) func() {
	var once sync.Once
	return func() {
		once.Do(fn)
	}
}

// OnceValue returns a function that invokes fn only once and returns the
// value returned by fn. The returned function may be called concurrently.
func OnceValue[T any](fn func() T) func() T {
	var (
		once sync.Once
		v    T
	)
	return func() T {
		once.Do(func() {
			v = fn()
		})
		return v
	}
}

// OnceValues returns a function that invokes fn only once and returns the
// values returned by fn. The returned function may be called concurrently.
func OnceValues[T1, T2 any](fn func() (T1, T2)) func() (T1, T2) {
	var (
		once sync.Once
		v1   T1
		v2   T2
	)
	return func() (T1, T2) {
		once.Do(func() {
			v1, v2 = fn()
		})
		return v1, v2
	}
}

// OnceFuncReflect returns a function that invokes fn only once and returns the
// value returned by fn. The returned function may be called concurrently.
func OnceFuncReflect[F any](fn F) F {
	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		panic(fmt.Errorf("sync: OnceFuncReflect called with non-function value (%T)", fn))
	}
	fnt := fnv.Type()
	if fnt.NumIn() != 0 {
		panic(fmt.Errorf("sync: OnceFuncReflect called with function with more than zero args (%T)", fn))
	}
	var (
		once    sync.Once
		results []reflect.Value
	)
	return reflect.MakeFunc(fnt, func([]reflect.Value) []reflect.Value {
		once.Do(func() {
			results = fnv.Call(nil)
		})
		return results
	}).Interface().(F)
}
