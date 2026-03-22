package core

import (
	"fmt"
	"testing"
)

func TestSimpleStringDecode(t *testing.T) {
	// Map of Network String -> Expected Parsed Value
	cases := map[string]string{
		"+OK\r\n": "OK",
	}

	for k, v := range cases {
		// Pass the network string into our Decode function
		value, _ := Decode([]byte(k))
		if v != value {
			t.Fail() // If the output doesn't match our expectation, fail the test!
		}
	}
}

func TestErrorDecode(t *testing.T) {
	cases := map[string]string{
		"-Error message\r\n": "Error message",
	}

	for k, v := range cases {
		value, _ := Decode([]byte(k))
		if v != value {
			t.Fail()
		}
	}
}

func TestInt64Decode(t *testing.T) {
	cases := map[string]int64{
		":0\r\n":    0,
		":1000\r\n": 1000,
	}

	for k, v := range cases {
		value, _ := Decode([]byte(k))
		if v != value {
			t.Fail()
		}
	}
}

func TestBulkStringDecode(t *testing.T) {
	cases := map[string]string{
		"$5\r\nhello\r\n": "hello",
		"$0\r\n\r\n":      "", // Tests an empty string!
	}

	for k, v := range cases {
		value, _ := Decode([]byte(k))
		if v != value {
			t.Fail()
		}
	}
}

func TestArrayDecode(t *testing.T) {
	// The value here is []interface{} because it's an array that can hold mixed data types
	cases := map[string][]interface{}{
		"*0\r\n":                                        {},
		"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n":          {"hello", "world"},
		"*3\r\n:1\r\n:2\r\n:3\r\n":                      {int64(1), int64(2), int64(3)},
		"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n": {int64(1), int64(2), int64(3), int64(4), "hello"},

		// This tests a nested array! An array inside an array.
		"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n": {[]interface{}{int64(1), int64(2), int64(3)}, []interface{}{"Hello", "World"}},
	}

	for k, v := range cases {
		value, _ := Decode([]byte(k))

		// Go Type Assertion: We must tell Go to treat 'value' as an array of interfaces
		array := value.([]interface{})

		// First, check if lengths match
		if len(array) != len(v) {
			t.Fail()
		}

		// Loop through every item in the array to ensure deep equality
		for i := range array {
			// fmt.Sprintf("%v") converts any data type to a string representation.
			// This is a clever trick Arpit uses to easily compare deeply nested arrays!
			if fmt.Sprintf("%v", v[i]) != fmt.Sprintf("%v", array[i]) {
				t.Fail()
			}
		}
	}
}
