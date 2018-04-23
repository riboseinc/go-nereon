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

	"github.com/hashicorp/hcl"
)

// parse configuration file
func (config *ConfigScheme) ParseConfigFile(configFpath string, cfgDict interface{}) error {
	// read configuration file
	cfg_data, err := ioutil.ReadFile(configFpath)
	if err != nil {
		return err
	}

	// get HCL tree
	cfg_tree, err := hcl.Parse(string(cfg_data))
	if err != nil {
		return err
	}

	// parse HCL configuration
	if err = hcl.DecodeObject(cfgDict, cfg_tree); err != nil {
		return err
	}

	return nil
}

// set configuration option
// func (config) SetCfgOption(key string, val interface{}) bool {
// 	var cfg_val interface{}

// 	seps := strings.Split(key, ".")

// 	for i:=0; i < len(seps); i++ {
// 		cfg_val = cfg.opts[seps[i]]
// 		if cfg_val == nil {
// 			return false
// 		}
// 	}

// 	return true
// }
