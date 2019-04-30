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
	"enode://c7c6a7e25dcca78d34c13dfc36b68ba5ecbd2fbc055e71107d80d3b6c1fd07ea2789bf195aac7c2686d4baca664d1aa5c04d4d7b4337f3a2213e2ab9b3a8ac16@121.141.157.230:40071",
	"enode://6a5e2c3e429e9f444259352fa1168e4cbe6afca5c72412f39b44ef01d7891151edaf40815d02c4005c0c9479e8f85cefc3317bc719678286c58d87a971c28604@121.141.157.228:40072",
	"enode://904605f4b510477ff59a76d2340ce35201e73c98ea81c1194d085da43c02f42d317bafaab3b7cb6d8afee459c7ca3b6533335ee4d132fc9b3cab5ef5d6f2639a@34.237.211.223:40071",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	"enode://c7c6a7e25dcca78d34c13dfc36b68ba5ecbd2fbc055e71107d80d3b6c1fd07ea2789bf195aac7c2686d4baca664d1aa5c04d4d7b4337f3a2213e2ab9b3a8ac16@121.141.157.230:40071",
	"enode://6a5e2c3e429e9f444259352fa1168e4cbe6afca5c72412f39b44ef01d7891151edaf40815d02c4005c0c9479e8f85cefc3317bc719678286c58d87a971c28604@121.141.157.228:40072",
	"enode://904605f4b510477ff59a76d2340ce35201e73c98ea81c1194d085da43c02f42d317bafaab3b7cb6d8afee459c7ca3b6533335ee4d132fc9b3cab5ef5d6f2639a@34.237.211.223:40071",
}


// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	"enode://c7c6a7e25dcca78d34c13dfc36b68ba5ecbd2fbc055e71107d80d3b6c1fd07ea2789bf195aac7c2686d4baca664d1aa5c04d4d7b4337f3a2213e2ab9b3a8ac16@121.141.157.230:40071",
	"enode://6a5e2c3e429e9f444259352fa1168e4cbe6afca5c72412f39b44ef01d7891151edaf40815d02c4005c0c9479e8f85cefc3317bc719678286c58d87a971c28604@121.141.157.228:40072",
	"enode://904605f4b510477ff59a76d2340ce35201e73c98ea81c1194d085da43c02f42d317bafaab3b7cb6d8afee459c7ca3b6533335ee4d132fc9b3cab5ef5d6f2639a@34.237.211.223:40071",
}
