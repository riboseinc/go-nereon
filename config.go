
package main

import (
	"fmt"
	"os"
	"errors"
	"io/ioutil"

	ast "github.com/hashicorp/hcl/hcl/ast"
	hclParser "github.com/hashicorp/hcl/hcl/parser"
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
	envName        string
	fullDesc       string
	children       []ConfigItemScheme
}

type ConfigScheme struct {
	configItems    []ConfigItemScheme
	configFpath    string

	cmdlineOpts    map[string]interface{}
	envOpts        map[string]interface{}
	cfgOpts        map[string]interface{}
}

// create new config scheme
func NewConfigScheme(configItems []ConfigItemScheme, configFpath string) *ConfigScheme {
	return &ConfigScheme {
		configItems:        configItems,
		configFpath:        configFpath,

		cmdlineOpts:        make(map[string]interface{}),
		envOpts:            make(map[string]interface{}),
		cfgOpts:            make(map[string]interface{}),
	}
}

// merge each option values
func (config *ConfigScheme) MergeOpts() (map[string]interface{}, error) {
	return nil, nil
}

// parse configuration
func (config *ConfigScheme) ParseConfig() (map[string]interface{}, error) {
	var err error

	// parse command line arguments
	err = config.ParseCmdLine()
	if err != nil {
		return nil, err
	}

	// parse configuration file
	if len(config.configFpath) == 0 {
		if config.cmdlineOpts["config"] != nil {
			config.configFpath = config.cmdlineOpts["config"].(string)
		} else if config.envOpts["config"] != nil {
			config.configFpath = config.envOpts["config"].(string)
		}
	}
	
	err = config.ParseConfigFile()
	if err != nil {
		return nil, err
	}

	// parse environment variable
	err = config.ParseEnvVars()
	if err != nil {
		return nil, err
	}

	return config.MergeOpts()
}

// check type of option value
func CheckOptValType(opt_type ConfigValueType, opt_val interface{}) bool {
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

	return matched
}

// parse command line arguments
func (config *ConfigScheme) ParseCmdLine() error {
	opts := make(map[string]interface{})

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
				opts[item.name] = true
				found = true
				break
			}

			// check for variable type
			i++
			if i == len(os.Args) {
				return errors.New(fmt.Sprintf("Missing argument for option '%s'", arg))
			}

			if matched := CheckOptValType(item.valType, os.Args[i]); !matched {
				return errors.New(fmt.Sprintf("Invalid argument type for option '%s'", arg))
			}

			// set value
			opts[item.name] = string(os.Args[i])
			found = true

			break
		}

		if !found {
			return errors.New(fmt.Sprintf("Invalid option '%s'", arg))
		}
	}

	config.cmdlineOpts = opts

	return nil
}

// parse configuration file
func (config *ConfigScheme) ParseConfigFile() error {
	var cfg_data []byte
	var err error

	var astFile *ast.File

	// check the configuration file is exist
	if len(config.configFpath) == 0 {
		return nil
	}

	// read configuration file
	if cfg_data, err = ioutil.ReadFile(config.configFpath); err != nil {
		return errors.New("Could not read configuration file")
	}

	// parse HCL configuration
	if astFile, err = hclParser.Parse(cfg_data); err != nil {
		return err
	}
	fmt.Println(astFile)

	return nil
}

// parse environment variables
func (config *ConfigScheme) ParseEnvVars() error {
	return nil
}

// print help
func (config *ConfigScheme) PrintHelp() {
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
