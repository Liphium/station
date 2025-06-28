package magic_util

import "log"

func AssertEq(actual interface{}, expectation interface{}) {
	if actual != expectation {
		log.Fatalf("got %s, expected %s", actual, expectation)
	}
}
