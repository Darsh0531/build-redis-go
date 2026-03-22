package core

import (
	"errors"
)

// Decode is the main entry point. It takes raw bytes and returns a Go data type (interface{}).
func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	// DecodeOne returns the value, how many bytes it read (delta), and an error.
	// We use "_" to ignore the delta here since we just want the final value.
	value, _, err := DecodeOne(data)
	return value, err
}

// DecodeOne checks the first character and routes it to the correct reading function.
func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	// data[0] looks at the very first byte (e.g., '+', '-', '*', etc.)
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}

	return nil, 0, nil
}

// +OK\r\n
func readSimpleString(data []byte) (string, int, error) {
	pos := 1 // Start at index 1 to skip the '+'

	// Keep moving 'pos' forward until we hit '\r' (Carriage Return)
	for ; data[pos] != '\r'; pos++ {
	}

	// Extract the bytes from index 1 to pos, and convert them to a String.
	// Return the string, pos + 2 (to skip \r and \n for the next read), and nil (no error)
	return string(data[1:pos]), pos + 2, nil
}

// -Error message\r\n
func readError(data []byte) (string, int, error) {
	// An error string is parsed exactly the same way as a simple string!
	return readSimpleString(data)
}

// :1000\r\n
func readInt64(data []byte) (int64, int, error) {
	pos := 1
	var value int64 = 0

	for ; data[pos] != '\r'; pos++ {
		// ASCII trick: data[pos] is a byte (e.g., '5' is byte 53).
		// byte '0' is 48. So 53 - 48 = 5.
		// value = value * 10 shifts digits left (e.g., 1 -> 10 -> 105).
		value = value*10 + int64(data[pos]-'0')
	}

	return value, pos + 2, nil
}

// Helper to read lengths (used by Arrays and Bulk Strings)
// Reads digits until '\r\n'
func readLength(data []byte) (int, int) {
	pos, length := 0, 0
	for pos = range data {
		b := data[pos]
		if !(b >= '0' && b <= '9') {
			// As soon as we hit a non-number (like \r), return length and pos + 2 (skip \r\n)
			return length, pos + 2
		}
		length = length*10 + int(b-'0')
	}
	return 0, 0
}

// $5\r\nhello\r\n
func readBulkString(data []byte) (string, int, error) {
	pos := 1

	// 1. Read the length (e.g., "5")
	len, delta := readLength(data[pos:])
	pos += delta

	// 2. Extract exactly 'len' characters.
	// data[pos : pos+len] creates a sub-array of exactly that word!
	// Return pos + len + 2 (skip the final \r\n)
	return string(data[pos:(pos + len)]), pos + len + 2, nil
}

// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
func readArray(data []byte) (interface{}, int, error) {
	pos := 1

	// 1. Read the number of elements in the array
	count, delta := readLength(data[pos:])
	pos += delta

	// 2. Create an empty array (slice) of size `count` to hold our values
	var elems []interface{} = make([]interface{}, count)

	// 3. Loop exactly `count` times.
	for i := range elems {
		// Recursively call DecodeOne! It will automatically figure out if the
		// next item is a string, int, etc., process it, and tell us how many bytes it consumed.
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}

		elems[i] = elem // Save the decoded element
		pos += delta    // Move our pointer forward by the amount of bytes consumed
	}

	return elems, pos, nil
}
