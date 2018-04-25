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
	"io/ioutil"
	"strings"
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl"
)

// parse configuration file
func (config *ConfigScheme) ConvertConfigToMap(configFpath string) error {
	// read configuration file
	cfg_data, err := ioutil.ReadFile(configFpath)
	if err != nil {
		return err
	}

	if err := hcl.Unmarshal(cfg_data, &config.cfgMap); err != nil {
		return err
	}

	return nil
}

// set configuration option
func (config *ConfigScheme) OverrideCfgOption(key string, val interface{}) {
	var i int
	var cfg_map []map[string]interface{}

	seps := strings.Split(key, ".")
	if len(seps) == 0 {
		return
	}

	cfg_val := config.cfgMap[seps[0]]
	for i=1; i < len(seps); i++ {
		if cfg_val == nil {
			break
		}

		// check type
		cfg_type := fmt.Sprintf("%v", reflect.TypeOf(cfg_val))
		if cfg_type != "[]map[string]interface {}" {
			break
		}

		cfg_map = cfg_val.([]map[string]interface{})
		if i == len(seps) - 1 {
			break
		}
		cfg_val = (cfg_map[0])[seps[i]]
	}
	cfg_map[0][seps[len(seps) - 1]] = val

	return
}
