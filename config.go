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
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

// type for option value
const (
	OPT_TYPE_STRING = "string"
	OPT_TYPE_BOOL   = "bool"
	OPT_TYPE_INT    = "int"
	OPT_TYPE_IPPORT = "ipport"
	OPT_TYPE_ARRAY  = "array"
	OPT_TYPE_PROTO  = "proto"
)

type CliKeyOpts struct {
	ShortKey          string                `hcl:"short"`
	FullKey           string                `hcl:"long"`
}

type HCLCliOptions struct {
	Name              string                `hcl:",key"`
	Type              string                `hcl:"type"`
	Switch            CliKeyOpts            `hcl:"switch"`
	Desc              CliKeyOpts            `hcl:"description"`
	EnvName           string                `hcl:"env"`
	CfgKey            string                `hcl:"config"`
	ShowHelper        bool                  `hcl:"helper"`
	OverrideCfgPath   bool                  `hcl:"override_cfg"`
	Value             interface{}
}

type HCLCliConfig struct {
	HclCliOpts     []HCLCliOptions          `hcl:"cmdline"`
}

type ConfigScheme struct {
	hclCliConfig   *HCLCliConfig
	configFpath    string

	envConfig      *EnvConfig
}

// create new config scheme
func NewConfigScheme() *ConfigScheme {
	return &ConfigScheme {
		hclCliConfig:        &HCLCliConfig{},
		envConfig:           NewEnvConfig(),
	}
}

// parse HCL options
func (config *ConfigScheme) ParseHCLOptions(hclOptFpath string) error {
	// read HCL configuration from file
	hclBytes, err := ioutil.ReadFile(hclOptFpath)
	if err != nil {
		return err
	}

	// decode HCL object
	hclTree, err := hcl.Parse(string(hclBytes))
	if err != nil {
		return err
	}

	if err := hcl.DecodeObject(&config.hclCliConfig, hclTree); err != nil {
		return err
	}

	return nil
}

// parse command line arguments
func (config *ConfigScheme) ParseCmdLine() error {
	for i:=1; i < len(os.Args); i++ {
		found := false
		arg := os.Args[i]

		for j:=0; j < len(config.hclCliConfig.HclCliOpts); j++ {
			opt := config.hclCliConfig.HclCliOpts[j]

			shortKey := "-" + opt.Switch.ShortKey
			fullKey := "--" + opt.Switch.FullKey

			// matches for short and long name
			if arg != shortKey && arg != fullKey {
				continue
			}

			// check whether argument needs value
			if opt.Type == OPT_TYPE_BOOL {
				found = true
				break
			}

			if matched := CheckOptValType(opt.Name, os.Args[i], opt.Type); !matched {
				return errors.New(fmt.Sprintf("Invalid value type '%v' for '%s' option", os.Args[i], opt.Name))
			}

			// check for variable type
			i++
			if i == len(os.Args) {
				return errors.New(fmt.Sprintf("Missing argument for option '%s'", arg))
			}

			found = true
			break
		}

		if !found {
			return errors.New(fmt.Sprintf("Invalid option '%s'", arg))
		}
	}

	return nil
}

// print help
func (config *ConfigScheme) PrintCmdLineHelp() {
	maxFullkeyLen := 0
	maxDescLen := 0
	for i:=0; i < len(config.hclCliConfig.HclCliOpts); i++ {
		opt := config.hclCliConfig.HclCliOpts[i]

		if len(opt.Switch.FullKey) > maxFullkeyLen {
			maxFullkeyLen = len(opt.Switch.FullKey)
		}

		if len(opt.Desc.ShortKey) > maxDescLen {
			maxDescLen = len(opt.Desc.ShortKey)
		}
	}
	maxFullkeyLen += 2
	maxDescLen += 2

	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	for i:=0; i < len(config.hclCliConfig.HclCliOpts); i++ {
		opt := config.hclCliConfig.HclCliOpts[i]

		paddingFullkey := FillBytesArray(maxFullkeyLen - len(opt.Switch.FullKey), ' ')
		paddingDesc := FillBytesArray(maxDescLen - len(opt.Desc.ShortKey), ' ')
		if opt.Type != OPT_TYPE_BOOL {
			fmt.Printf("  -%v|--%v%v<%v>%v: %v\n",
				opt.Switch.ShortKey,
				opt.Switch.FullKey, string(paddingFullkey),
				opt.Desc.ShortKey, string(paddingDesc),
				opt.Desc.FullKey)
		} else {
			fmt.Printf("  -%v|--%v%v%v  : %v\n",
				opt.Switch.ShortKey,
				opt.Switch.FullKey, string(paddingFullkey),
				string(paddingDesc),
				opt.Desc.FullKey)
		}
	}
}
