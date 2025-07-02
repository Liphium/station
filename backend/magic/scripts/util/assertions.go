package magic_util

import (
	"maps"
	"reflect"
	"runtime/debug"
	"testing"

	"github.com/Liphium/station/neogate"
)

func AssertEq[T comparable](t *testing.T, actual T, expectation T) {
	assertFunc(t, actual, expectation, func(i1, i2 T) bool {
		return i1 == i2
	})
}

func AssertDeepEq(t *testing.T, actual interface{}, expectation interface{}) {
	assertFunc(t, actual, expectation, func(i1, i2 interface{}) bool {
		return reflect.DeepEqual(i1, i2)
	})
}

func AssertEventsEq(t *testing.T, actual, expectation neogate.Event) {
	assertFunc(t, actual, expectation, func(i1, i2 neogate.Event) bool {
		return i1.Name == i2.Name && maps.Equal(i1.Data, i2.Data)
	})
}

func assertFunc[T any](t *testing.T, actual T, expectation T, comparer func(T, T) bool) {
	if !comparer(actual, expectation) {
		debug.PrintStack()
		t.Fatalf("Assertion failed: got %v, expected %v", actual, expectation)
	}
}
