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
	"enode://332569f6bbb365006eecffa5578a14863e224577355295bf8e3dcef8e93501fb3b94f1f0a0f8ce29cb7dfc0ad42174db8e4642835b6f10330ee51257138c336b@121.141.157.230:40404", //MEGA
	"enode://248e0835b6fc0449622ef18d16a123485bb1d66b40c85b8b9a1edb5a86ac51feb78564ce84e5edd16472df0738a60ffd3a8e1353a9ab8d21fdcde880583a6d90@121.141.157.228:40404", //GIGA
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	"enode://332569f6bbb365006eecffa5578a14863e224577355295bf8e3dcef8e93501fb3b94f1f0a0f8ce29cb7dfc0ad42174db8e4642835b6f10330ee51257138c336b@121.141.157.230:40404", //MEGA
	"enode://248e0835b6fc0449622ef18d16a123485bb1d66b40c85b8b9a1edb5a86ac51feb78564ce84e5edd16472df0738a60ffd3a8e1353a9ab8d21fdcde880583a6d90@121.141.157.228:40404", //GIGA
}


// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	//"enode://06051a5573c81934c9554ef2898eb13b33a34b94cf36b202b69fde139ca17a85051979867720d4bdae4323d4943ddf9aeeb6643633aa656e0be843659795007a@35.177.226.168:30303",
	//"enode://0cc5f5ffb5d9098c8b8c62325f3797f56509bff942704687b6530992ac706e2cb946b90a34f1f19548cd3c7baccbcaea354531e5983c7d1bc0dee16ce4b6440b@40.118.3.223:30304",
	//"enode://1c7a64d76c0334b0418c004af2f67c50e36a3be60b5e4790bdac0439d21603469a85fad36f2473c9a80eb043ae60936df905fa28f1ff614c3e5dc34f15dcd2dc@40.118.3.223:30306",
	//"enode://85c85d7143ae8bb96924f2b54f1b3e70d8c4d367af305325d30a61385a432f247d2c75c45c6b4a60335060d072d7f5b35dd1d4c45f76941f62a4f83b6e75daaf@40.118.3.223:30307",

	"enode://332569f6bbb365006eecffa5578a14863e224577355295bf8e3dcef8e93501fb3b94f1f0a0f8ce29cb7dfc0ad42174db8e4642835b6f10330ee51257138c336b@121.141.157.230:40404", //MEGA
	"enode://248e0835b6fc0449622ef18d16a123485bb1d66b40c85b8b9a1edb5a86ac51feb78564ce84e5edd16472df0738a60ffd3a8e1353a9ab8d21fdcde880583a6d90@121.141.157.228:40404", //GIGA
	//"enode://fd3c3f53a5fb15c8ef981c1117849efeb7259fb14c5c9345c62b4e812ada653bfdada15e9492863798c3b0fc3f2925e0800da5a458e7308a788287d541d79df8@192.168.0.160:30310",
}
