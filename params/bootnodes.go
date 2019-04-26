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

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Berith network.
var MainnetBootnodes = []string{
	// Berith Foundation Go Bootnodes
	"enode://f9781cd34f75032cce7d8f5a538a29ffcc02af09713b36a00e1a44551ac3530859b112493973aad54774f92c62183e2837352a8e5933793f91e11a374992aaae@34.237.211.223:49999",
	//"enode://7fdd9b52c173a3a10d4a2718a3b3712fd2c69cd622555706f22ed45c1ceff6e39cdb39acddf04abbe3e881986861e84971cc57601381d708e0c3ad9af397de2b@34.237.211.223:49999", //MEGA
	//"enode://248e0835b6fc0449622ef18d16a123485bb1d66b40c85b8b9a1edb5a86ac51feb78564ce84e5edd16472df0738a60ffd3a8e1353a9ab8d21fdcde880583a6d90@121.141.157.228:40404", //GIGA
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	"enode://f9781cd34f75032cce7d8f5a538a29ffcc02af09713b36a00e1a44551ac3530859b112493973aad54774f92c62183e2837352a8e5933793f91e11a374992aaae@34.237.211.223:49999",
	//"enode://7fdd9b52c173a3a10d4a2718a3b3712fd2c69cd622555706f22ed45c1ceff6e39cdb39acddf04abbe3e881986861e84971cc57601381d708e0c3ad9af397de2b@34.237.211.223:49999", //MEGA
	//"enode://248e0835b6fc0449622ef18d16a123485bb1d66b40c85b8b9a1edb5a86ac51feb78564ce84e5edd16472df0738a60ffd3a8e1353a9ab8d21fdcde880583a6d90@121.141.157.228:40404", //GIGA
}


// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	"enode://f9781cd34f75032cce7d8f5a538a29ffcc02af09713b36a00e1a44551ac3530859b112493973aad54774f92c62183e2837352a8e5933793f91e11a374992aaae@34.237.211.223:49999",

	//"enode://7fdd9b52c173a3a10d4a2718a3b3712fd2c69cd622555706f22ed45c1ceff6e39cdb39acddf04abbe3e881986861e84971cc57601381d708e0c3ad9af397de2b@34.237.211.223:49999", //MEGA
	//"enode://248e0835b6fc0449622ef18d16a123485bb1d66b40c85b8b9a1edb5a86ac51feb78564ce84e5edd16472df0738a60ffd3a8e1353a9ab8d21fdcde880583a6d90@121.141.157.228:40404", //GIGA
	//"enode://fd3c3f53a5fb15c8ef981c1117849efeb7259fb14c5c9345c62b4e812ada653bfdada15e9492863798c3b0fc3f2925e0800da5a458e7308a788287d541d79df8@192.168.0.160:30310",
}
