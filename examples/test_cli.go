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

	".."
)

type globalConfig struct {
	LogDir         string             `hcl:"log_directory"`
	LogLevel       int                `hcl:"log_level"`

	ListenAddr     string             `hcl:"listen_address"`
}

type cfgFilesConfig struct {
	Name           string             `hcl:",key"`
	Files          []string           `hcl:"files"`
}

type cmdConfig struct {
	Name           string             `hcl:",key"`
	Uid            int                `hcl:"uid"`
	Interval       int                `hcl:"interval"`
	Cmds           []string           `hcl:"cmds"`
}

type RsyslogConfig struct {
	ListenAddr     string             `hcl:"listen_address"`
	Protocol       string             `hcl:"protocol"`
}

type testConfig struct {
	Globals        globalConfig       `hcl:"global"`
	ConfigFiles    []cfgFilesConfig   `hcl:"config-files"`
	Commands       []cmdConfig        `hcl:"commands"`
	Rsyslog        RsyslogConfig      `hcl:"rsyslog"`
}

func main() {
	testCfg := &testConfig{}
	config := mconfig.NewConfigScheme()

	// parse HCL options
	if err := config.ParseConfig("options.example", "config.example", &testCfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(testCfg)
}
