
package main

import (
	"os"
	"fmt"
)

var opts = []ConfigItemScheme {
	{
		"config",
		"The configuration file path",
		CONF_VAL_TYPE_STRING,
		&CommandLineSwitch{'c', "config"},
		"",
		"Specify the configuration file path",
		nil,
	},

	{
		"log-dir",
		"The logging directory",
		CONF_VAL_TYPE_STRING,
		&CommandLineSwitch{'l', "log-dir"},
		"log-settings.directory",
		"Specify the logging directory",
		nil,
	},
}

func main() {
	var err error

	config := NewConfigScheme(opts, "examples/cfg/config.example")

	// parsing command line and configuration file
	opts := make(map[string]interface{})
	if opts, err = config.ParseConfig(); err != nil {
		fmt.Println("Parsing configuration has failed.")
		os.Exit(1)
	}

	// print options
	fmt.Println(opts)
}
