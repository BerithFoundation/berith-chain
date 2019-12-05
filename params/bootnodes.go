// Copyright 2015 The go-Berith Authors
// This file is part of the go-Berith library.
//
// The go-Berith library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-Berith library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-Berith library. If not, see <http://www.gnu.org/licenses/>.

/*
[BERITH]
부트노드 에 대한 정보를 적는 곳
*/

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Berith network.
var MainnetBootnodes = []string{
	// Berith Foundation Go Bootnodes
	"enode://b51be77505e2391e0e2ad13a5153607ab5fd19b89a513d6ba5bc793aff6643b3b06dc33b05245b417f5dc5fcbd6a755b679c5730059e15346ddb84c96adbe388@15.164.130.81:40404",
	"enode://6fdc64e4d2ad8c973a8606e4d1947b734a437d8f08f2c3e321c71f791b26ad27e38e0b9700900585a489373f38325f9a083b7e5a4520d9564c41ba9a44ec997e@13.124.140.180:40404",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	"enode://8142df22d4ca164db41f0cb6a66f439332ec8d9a799dbd1a5bbc990ca585fbc283dc949edb8a0962942accb0c361a1171dc16dd46abefd17b7dc8ef19e6894e8@13.209.204.31:55555",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
