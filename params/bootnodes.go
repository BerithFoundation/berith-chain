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
Information about boot nodes
*/

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Berith network.
var MainnetBootnodes = []string{
	// Berith Foundation Go Bootnodes
	"enode://ea9b7c833a522780cb50dbb5f6e44c8d475ce8dedda44cb555e59994a5f89288908ebb288cfec9962c7321dee311a2a9bbfbadda78b1b3ef6dbcb33aea063e21@13.124.140.180:40404",
	"enode://2c7f9c316e460f98516e27ecd4bcb3e6772d2ae50db7ed11f96b4cb973aaca51b21cb485815d9f627c607e9def084c6e183cd2c12ec9dcc22fd9af198b6d34d3@15.164.130.81:40404",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	"enode://088702216a04cd6bfa865e79392c998735961fc86670a39842735941beac25f8b78f0d7a72fe542d10622d45dddf40902c1737c6574c080d2aae655dc785ed92@15.164.124.7:30301",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
