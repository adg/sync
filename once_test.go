package sync

import (
	"errors"
	"reflect"
	"testing"
)

var oops = errors.New("oops")

func TestOnceFuncReflect(t *testing.T) {
	var calls int
	check := func() {
		if calls > 0 {
			panic("function called more than once")
		}
		calls++
	}
	for _, c := range []struct {
		f    any
		want []any
	}{
		{func() { check() }, nil},
		{func() int { check(); return 42 }, []any{42}},
		{func() (int, error) { check(); return 42, oops }, []any{42, oops}},
	} {
		t.Run(reflect.TypeOf(c.f).String(), func(t *testing.T) {
			calls = 0
			fn := OnceFuncReflect(c.f)
			for j := 0; j < 2; j++ {
				got := reflect.ValueOf(fn).Call(nil)
				if len(got) != len(c.want) {
					t.Fatalf("got %d results, want %d", len(got), len(c.want))
				}
				for i, v := range got {
					if g, w := v.Interface(), c.want[i]; g != w {
						t.Fatalf("result %d is %v, want %v", i, g, w)
					}
				}
			}
		})
	}
}

func BenchmarkOnceFunc(b *testing.B) {
	b.Run("Func", func(b *testing.B) {
		fn := OnceFunc(func() {})
		for i := 0; i < b.N; i++ {
			fn()
		}
	})
	b.Run("Value", func(b *testing.B) {
		fn := OnceValue(func() int { return 42 })
		for i := 0; i < b.N; i++ {
			if v := fn(); v != 42 {
				b.Fatalf("got %v, want 42", v)
			}
		}
	})
	b.Run("Values", func(b *testing.B) {
		fn := OnceValues(func() (int, error) { return 42, oops })
		for i := 0; i < b.N; i++ {
			if v, err := fn(); v != 42 || err != oops {
				b.Fatalf("got %v, %v, want 42, oops", v, err)
			}
		}
	})
}

func BenchmarkOnceFuncReflect(b *testing.B) {
	b.Run("Func", func(b *testing.B) {
		fn := OnceFuncReflect(func() {})
		for i := 0; i < b.N; i++ {
			fn()
		}
	})
	b.Run("Value", func(b *testing.B) {
		fn := OnceFuncReflect(func() int { return 42 })
		for i := 0; i < b.N; i++ {
			if v := fn(); v != 42 {
				b.Fatalf("got %v, want 42", v)
			}
		}
	})
	b.Run("Values", func(b *testing.B) {
		fn := OnceFuncReflect(func() (int, error) { return 42, oops })
		for i := 0; i < b.N; i++ {
			if v, err := fn(); v != 42 || err != oops {
				b.Fatalf("got %v, %v, want 42, oops", v, err)
			}
		}
	})
}
