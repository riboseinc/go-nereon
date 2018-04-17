/*
 * Copyright (c) 2018, [Ribose Inc](https://www.ribose.com).
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
	"os"
	"errors"
)

type ConfigValueType uint32
const (
	CONF_VAL_TYPE_INT ConfigValueType = iota
	CONF_VAL_TYPE_STRING
	CONF_VAL_TYPE_FLOAT
	CONF_VAL_TYPE_BOOLEAN
	CONF_VAL_TYPE_IPPORT
	CONF_VAL_TYPE_HOSTPORT
	CONF_VAL_TYPE_IPADDR
	CONF_VAL_TYPE_ARRAY
)

type CommandLineSwitch struct {
	shortKey       byte
	fullKey        string
}

type ConfigItemScheme struct {
	name           string
	shortDesc      string
	valType        ConfigValueType
	command        *CommandLineSwitch
	cfgName        string
	envName        string
	fullDesc       string
	children       []ConfigItemScheme
}

type ConfigScheme struct {
	configItems    []ConfigItemScheme
	configFpath    string

	cfg            *HCLConfig
	ev             *EnvConfig

	opts           map[string]interface{}
}

// create new config scheme
func NewConfigScheme(configItems []ConfigItemScheme, configFpath string) *ConfigScheme {
	return &ConfigScheme {
		configItems:        configItems,
		configFpath:        configFpath,

		cfg:                NewHCLConfig(),
		ev:                 NewEnvConfig(),

		opts:               make(map[string]interface{}),
	}
}

// parse configuration
func (config *ConfigScheme) ParseConfig() error {
	var err error

	// parse configuration file
	if err := config.cfg.ParseConfigFile(config.configFpath); err != nil {
		return err
	}

	// parse configuration items
	if err := config.ev.ParseEnv(config.configItems, config.cfg.opts); err != nil {
		return err
	}

	// parse command line arguments
	if err = config.ParseCmdLine(); err != nil {
		return err
	}

	return nil
}

// parse command line arguments
func (config *ConfigScheme) ParseCmdLine() error {
	for i:=1; i < len(os.Args); i++ {
		found := false
		arg := os.Args[i]

		for j:=0; j < len(config.configItems); j++ {
			item := config.configItems[j]

			shortKey := "-" + string(item.command.shortKey)
			fullKey := "--" + item.command.fullKey

			// matches for short and long name
			if arg != shortKey && arg != fullKey {
				continue
			}

			// check whether argument needs value
			if item.valType == CONF_VAL_TYPE_BOOLEAN {
				config.opts[item.name] = true
				found = true
				break
			}

			// check for variable type
			i++
			if i == len(os.Args) {
				return errors.New(fmt.Sprintf("Missing argument for option '%s'", arg))
			}

			if err := CheckOptValType(item.name, os.Args[i], item.valType); err != nil {
				return err
			}

			// check if the value was set by configuration file
			if cfg_exist := config.cfg.SetCfgOption(item.cfgName, os.Args[i]); cfg_exist {
				found = true
				continue
			}

			// check if the value was set by environment
			if config.ev.opts[item.name] != nil {
				fmt.Printf("Found val '%v' for name '%s' from environment\n", config.cfg.opts[item.name], item.name)
				config.ev.opts[item.name] = os.Args[i]
				found = true
				continue
			}

			fmt.Printf("name: %s, val: %v\n", item.name, os.Args[i])

			config.opts[item.name] = os.Args[i]
			found = true

			break
		}

		if !found {
			return errors.New(fmt.Sprintf("Invalid option '%s'", arg))
		}
	}

	for k, v := range config.cfg.opts {
		config.opts[k] = v
	}

	for k, v := range config.ev.opts {
		config.opts[k] = v
	}

	return nil
}

// print help
func (config *ConfigScheme) PrintCmdLineHelp() {
	maxFullkeyLen := 0
	maxDescLen := 0
	for i:=0; i < len(config.configItems); i++ {
		item := config.configItems[i]

		if len(item.command.fullKey) > maxFullkeyLen {
			maxFullkeyLen = len(item.command.fullKey)
		}

		if len(item.shortDesc) > maxDescLen {
			maxDescLen = len(item.shortDesc)
		}
	}
	maxFullkeyLen += 2
	maxDescLen += 2

	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	for i:=0; i < len(config.configItems); i++ {
		item := config.configItems[i]

		paddingFullkey := FillBytesArray(maxFullkeyLen - len(item.command.fullKey), ' ')
		paddingDesc := FillBytesArray(maxDescLen - len(item.shortDesc), ' ')
		if item.valType != CONF_VAL_TYPE_BOOLEAN {
			fmt.Printf("  -%v|--%v%v<%v>%v: %v\n",
				string(item.command.shortKey),
				item.command.fullKey, string(paddingFullkey),
				item.shortDesc, string(paddingDesc),
				item.fullDesc)
		} else {
			fmt.Printf("  -%v|--%v%v%v  : %v\n",
				string(item.command.shortKey),
				item.command.fullKey, string(paddingFullkey),
				string(paddingDesc),
				item.fullDesc)
		}
	}
}
