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

package main

import (
	"os"
	"fmt"

	"github.com/riboseinc/go-multiconfig"
)

var opts = []ConfigItemScheme {
	{
		"config",
		"The configuration file path",
		CONF_VAL_TYPE_STRING,
		&CommandLineSwitch{'c', "config"},
		"", "",
		"Specify the configuration file path",
		nil,
	},

	{
		"log-dir",
		"The logging directory",
		CONF_VAL_TYPE_STRING,
		&CommandLineSwitch{'d', "log-dir"},
		"log-settings.directory", "",
		"Specify the logging directory",
		nil,
	},

	{
		"listen",
		"Listening address for incoming events",
		CONF_VAL_TYPE_IPADDR,
		&CommandLineSwitch{'l', "listen"},
		"listen.address", "",
		"Specify the listenning address",
		nil,
	},
}

func main() {
	var err error

	config := mconfig.NewConfigScheme(opts, "examples/cfg/config.example")

	// parsing command line and configuration file
	if err = config.ParseConfig(); err != nil {
		fmt.Println(err)

		// print command line helper
		config.PrintCmdLineHelp()

		os.Exit(1)
	}

	for i:=0; i < len(opts); i++ {
		msg := fmt.Sprintf("The value for option '%s' is '%v'", opts[i].name, config.opts[opts[i].name])
		fmt.Println(msg)
	}
}
