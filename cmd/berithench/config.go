// Copyright 2019 The berith Authors
// This file is part of berith.
//
// berith is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// berith is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with berith. If not, see <http://www.gnu.org/licenses/>.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/BerithFoundation/berith-chain/cmd/utils"
	"github.com/naoina/toml"
	cli "gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"
	"unicode"
)

type berithenchConfig struct {
	ChainID          int64    `json:"chainId"`          // chain id
	Nodes            []string `json:"nodes"`            // endpoints of nodes
	Keystore         string   `json:"keystore"`         // keystore path
	Addresses        []string `json:"addresses"`        // from addresses in tx
	Password         string   `json:"password"`         // password file path
	Duration         string   `json:"duration"`         // time duration for test execution
	TxCount          uint64   `json:"txCount"`          // tx count of test execution
	InitDelay        uint64   `json:"initDelay"`        // initial sleep before testing
	TxInterval       uint64   `json:"txInterval"`       // interval between send transactions
	OutputPath       string   `json:"outputPath"`       // path of results
	EnableCpuProfile bool     `json:"enableCpuProfile"` // enable cpu profile
	EnableGoTrace    bool     `json:"enableGoTrace"`    // enable go trace
}

var (
	defaultConfig = berithenchConfig{
		Keystore:         "",
		Password:         "",
		Duration:         "",
		TxCount:          0,
		TxInterval:       10,
		InitDelay:        0,
		OutputPath:       getDefaultWorkingDir(),
		EnableCpuProfile: false,
		EnableGoTrace:    false,
	}
	tomlSettings = toml.Config{
		NormFieldName: func(rt reflect.Type, key string) string {
			return key
		},
		FieldToKey: func(rt reflect.Type, field string) string {
			return field
		},
		MissingField: func(rt reflect.Type, field string) error {
			link := ""
			if unicode.IsUpper(rune(rt.Name()[0])) && rt.PkgPath() != "main" {
				link = fmt.Sprintf(", see https://godoc.org/%s#%s for available fields", rt.PkgPath(), rt.Name())
			}
			return fmt.Errorf("field '%s' is not defined in %s%s", field, rt.String(), link)
		},
	}
)

// parseConfig parse bbench config from cli context.
func parseConfig(ctx *cli.Context) (*berithenchConfig, error) {
	cfg := defaultConfig

	// load config file
	if file := ctx.String(ConfigFileFlag.Name); file != "" {
		if err := loadConfig(file, &cfg); err != nil {
			utils.Fatalf("%v", err)
		}
	}

	// apply flags
	var (
		chainID          = ctx.Int64(ChainIDFlag.Name)
		nodes            = ctx.String(NodesFlag.Name)
		keystore         = ctx.String(KeystoreFlag.Name)
		addrs            = ctx.String(AddressesFlag.Name)
		password         = ctx.String(PasswordFlag.Name)
		duration         = ctx.String(DurationFlag.Name)
		txCount          = ctx.Uint64(TxCountFlag.Name)
		txInterval       = ctx.Uint64(TxIntervalFlag.Name)
		initDelay        = ctx.Uint64(InitDelay.Name)
		output           = ctx.String(OutputPath.Name)
		enableCpuProfile = ctx.IsSet(EnableCpuProfile.Name)
		enableGoTrace    = ctx.IsSet(EnableGoTrace.Name)
	)
	// FIXME : more efficient compare
	if chainID != 0 {
		cfg.ChainID = chainID
	}
	if nodes != "" {
		cfg.Nodes = strings.Split(nodes, ",")
	}
	if keystore != "" {
		cfg.Keystore = keystore
	}
	if addrs != "" {
		cfg.Addresses = strings.Split(addrs, ",")
	}
	if password != "" {
		cfg.Password = password
	}
	if duration != "" {
		cfg.Duration = duration
	}
	if txCount > 0 {
		cfg.TxCount = txCount
	}
	if txInterval > 0 {
		cfg.TxInterval = txInterval
	}
	if initDelay > 0 {
		cfg.InitDelay = initDelay
	}
	if output != "" {
		cfg.OutputPath = output
	}
	if enableCpuProfile {
		cfg.EnableCpuProfile = enableCpuProfile
	}
	if enableGoTrace {
		cfg.EnableGoTrace = enableGoTrace
	}
	return &cfg, nil
}

// loadConfig load bbench config from toml file.
func loadConfig(file string, cfg *berithenchConfig) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	return err
}

// makePasswordList reads password lines from the file
func makePasswordList(config *berithenchConfig) []string {
	if config.Password == "" {
		return nil
	}
	content, err := ioutil.ReadFile(config.Password)
	if err != nil {
		utils.Fatalf("failed to read password file: %v", err)
	}
	var ret []string
	lines := strings.Split(string(content), "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
		if lines[i] != "" {
			ret = append(ret, lines[i])
		}
	}
	return ret
}

// parseDuration parse duration string to time with HH:mm:DD format
func parseDuration(config *berithenchConfig) time.Duration {
	if config.Duration == "" {
		return time.Duration(0)
	}
	t, err := time.Parse("15:04:05", config.Duration)
	if err != nil {
		utils.Fatalf("invalid duration. needs HH:mm:DD but %s.\n%v", config.Duration, err)
	}

	d := 0
	d += t.Hour() * int(time.Hour)
	d += t.Minute() * int(time.Minute)
	d += t.Second() * int(time.Second)
	return time.Duration(d)
}

// getWorkingDir returns user's home dir or temp dir if occur error
func getDefaultWorkingDir() string {
	home, err := os.UserHomeDir()
	if err == nil {
		return home
	}
	return os.TempDir()
}
