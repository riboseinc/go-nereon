/*
 * Copyright (c) 2017, [Ribose Inc](https://www.ribose.com).
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * ``AS IS'' AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package mconfig

import (
	"fmt"
	"strconv"
	"errors"
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

// check type of option value
func CheckOptValType(opt_name string, opt_val interface{}, opt_type ConfigValueType,) error {
	matched := false

	switch opt_type {
	case CONF_VAL_TYPE_INT:
		matched = ParseOptInt(opt_val)
	case CONF_VAL_TYPE_STRING:
		matched = ParseOptString(opt_val)
	case CONF_VAL_TYPE_ARRAY:
		matched = ParseOptArray(opt_val)
	case CONF_VAL_TYPE_IPADDR:
		matched = ParseOptAddrPair(opt_val)
	}

	if !matched {
		return errors.New(fmt.Sprintf("Invalid value type '%v' for '%s' option", opt_val, opt_name))
	}

	return nil
}

// fill array with specified character
func FillBytesArray(len int, ch byte) []byte {
	bytes_arr := make([]byte, len)
	for i := range bytes_arr {
		bytes_arr[i] = ch
	}

	return bytes_arr
}
