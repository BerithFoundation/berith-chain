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
	"fmt"
	"os"

	"github.com/BerithFoundation/berith-chain/cmd/utils"
	cli "gopkg.in/urfave/cli.v1"
)

var gitCommit = "" // Git SHA1 commit hash of the release (set via linker flags)
var app *cli.App

// Commonly used flags in cli.
var (
	ChainIDFlag = cli.Int64Flag{
		Name:  "chainid",
		Usage: "network chain id",
	}
	NodesFlag = cli.StringFlag{
		Name:  "nodes",
		Usage: "Comma separated list of nodes to send tx",
	}
	ConfigFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "Directory of config toml file",
	}
	KeystoreFlag = cli.StringFlag{
		Name:  "keystore",
		Usage: "Directory of keystore file",
	}
	PasswordFlag = cli.StringFlag{
		Name:  "password",
		Usage: "Password file path",
	}
	TxCountFlag = cli.Uint64Flag{
		Name:  "txcount",
		Usage: "How many test runs will be executed",
	}
	TxIntervalFlag = cli.Uint64Flag{
		Name:  "txinterval",
		Usage: "Interval between transactions [ms]",
	}
)

func init() {
	app = utils.NewAppWithHelpTemplate(gitCommit, "an Berith test tools", false)
	app.Commands = []cli.Command{
		ExecuteCommand,
		TpsCommand,
		AgentCommand,
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
