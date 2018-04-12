
package main

import (
	"strconv"
	"net"
)

// parse option for int type
func ParseOptInt(opt_val interface{}) bool {
	var str string
	var found bool

	if str, found = opt_val.(string); !found {
		return false
	}

	if _, err := strconv.Atoi(str); err != nil {
		return false
	}

	return true
}

// parse option for string type
func ParseOptString(opt_val interface{}) bool {
	_, found := opt_val.(string)
	return found
}

// parse option for array type
func ParseOptArray(opt_val interface{}) bool {
	str, found := opt_val.(string)
	if !found {
		return false
	}

	arr := []string{str}
	if len(arr) == 0 {
		return false
	}

	return true
}

// parse option for address type
func ParseOptAddrPair(opt_val interface{}) bool {
	str, found := opt_val.(string)
	if !found {
		return false
	}

	_, _, err := net.SplitHostPort(str)
	if err != nil {
		return false
	}

	return true
}

// fill array with specified character
func FillBytesArray(len int, ch byte) []byte {
	bytes_arr := make([]byte, len)
	for i := range bytes_arr {
		bytes_arr[i] = ch
	}

	return bytes_arr
}
