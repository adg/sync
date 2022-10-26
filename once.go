package sync

import "sync"

// OnceFunc returns a function that invokes fn only once and returns the values
// returned by fn. The returned function may be called concurrently.
func OnceFunc[T any](fn func() (T, error)) func() (T, error) {
	var (
		once  sync.Once
		value T
		err   error
	)
	return func() (T, error) {
		once.Do(func() {
			value, err = fn()
		})
		return value, err
	}
}

type OnceValue[T any] struct {
	once sync.Once
	v    T
}

func (o *OnceValue[T]) Do(fn func() T) T {
	o.once.Do(func() {
		o.v = fn()
	})
	return o.v
}

type OnceValueErr[T any] struct {
	once sync.Once
	v    T
	err  error
}

func (o *OnceValueErr[T]) Do(fn func() (T, error)) (T, error) {
	o.once.Do(func() {
		o.v, o.err = fn()
	})
	return o.v, o.err
}
