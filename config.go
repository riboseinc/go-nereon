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
	"encoding/json"

	"github.com/hashicorp/hcl"
	jsonParser "github.com/hashicorp/hcl/json/parser"
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
	hclCliConfig    *HCLCliConfig
	printHelpMsg    bool
	overrideCfgPath string

	envConfig       *EnvConfig

	cfgMap          map[string]interface{}
}

// create new config scheme
func NewConfigScheme() *ConfigScheme {
	return &ConfigScheme {
		hclCliConfig:        &HCLCliConfig{},
		printHelpMsg:        false,
		overrideCfgPath:     "",

		envConfig:           NewEnvConfig(),
		cfgMap:              make(map[string]interface{}),
	}
}

// parse the configuration file
func (config *ConfigScheme) ParseConfig(hclOpts string, cfgFpath string, cfgDict interface{}) error {
	var cfg string

	// parse command line at first
	if err := config.ParseHCLOptions(hclOpts); err != nil {
		fmt.Println("Parsing HCL options has failed.")
		return err
	}

	if config.overrideCfgPath != "" {
		cfg = config.overrideCfgPath
	} else {
		cfg = cfgFpath
	}

	// convert configuration to map
	if err := config.ConvertConfigToMap(cfg); err != nil {
		fmt.Println("Converting configuration to map has failed.")
		return err
	}

	// merge config
	config.MergeConfig()

	// marshal map to JSON
	json, err := json.MarshalIndent(config.cfgMap, "", " ")
	if err != nil {
		return err
	}

	ast, err := jsonParser.Parse(json)
	if err != nil {
		return err
	}

	if err := hcl.DecodeObject(cfgDict, ast); err != nil {
		return err
	}

	return nil
}

func (config *ConfigScheme) MergeConfig() {
	for i:=0; i < len(config.hclCliConfig.HclCliOpts); i++ {
		opt := config.hclCliConfig.HclCliOpts[i]

		if opt.CfgKey == "" || opt.Value == nil {
			continue
		}

		config.OverrideCfgOption(opt.CfgKey, opt.Value)
	}
}

// parse HCL options
func (config *ConfigScheme) ParseHCLOptions(hclOpts string) error {
	// decode HCL object
	hclTree, err := hcl.Parse(hclOpts)
	if err != nil {
		return err
	}

	if err := hcl.DecodeObject(&config.hclCliConfig, hclTree); err != nil {
		return err
	}

	return config.ParseCmdLine()
}

// parse command line arguments
func (config *ConfigScheme) ParseCmdLine() error {
	for i:=1; i < len(os.Args); i++ {
		found := false
		arg := os.Args[i]

		for j:=0; j < len(config.hclCliConfig.HclCliOpts); j++ {
			opt := &config.hclCliConfig.HclCliOpts[j]

			shortKey := "-" + opt.Switch.ShortKey
			fullKey := "--" + opt.Switch.FullKey

			// matches for short and long name
			if arg != shortKey && arg != fullKey {
				continue
			}

			// check whether argument needs value
			if opt.Type == OPT_TYPE_BOOL {
				if opt.ShowHelper == true {
					config.PrintCmdLineHelp()
					os.Exit(0)
				} else {
					opt.Value = true
				}

				found = true
				break
			}

			// check for variable type
			i++
			if i == len(os.Args) {
				return errors.New(fmt.Sprintf("Missing argument for option '%s'", arg))
			}

			if matched := CheckOptValType(opt.Name, os.Args[i], opt.Type); !matched {
				return errors.New(fmt.Sprintf("Invalid type '%v' for '%s' option value '%s'", opt.Type, opt.Name, os.Args[i]))
			}

			if opt.OverrideCfgPath == true {
				config.overrideCfgPath = os.Args[i]
			} else {
				opt.Value = os.Args[i]
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
