// Copyright 2017 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"fmt"
	"os"
)

// makeWizard creates and returns a new puppeth wizard.
func makeWizard() *wizard {
	return &wizard{
		conf: config{
			Servers: make(map[string][]byte),
		},
		in:       bufio.NewReader(os.Stdin),
	}
}

// run displays some useful infos to the user, starting on the journey of
// setting up a new or managing an existing Ethereum private network.
func (w *wizard)run() {
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println("| Welcome to Berith Genesis Wizard							|")
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println()

	w.makeGenesis()
}
