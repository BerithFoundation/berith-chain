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
	"enode://58e60ac860f35cf900b1fe92d51e419652651cff0cd7de085e5f415762a99b26c0c334d9a4a219ffd3f5032d42e5a0f446e1fca5b15a484f17442123f68d4aec@127.0.0.1:40404",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	"enode://26a8d8bd85d676b8ce8fb1afd5e75492224f10dbe911c493ecc275eb4eddab273a0e73f93c78eb2a3375a04be34b5afb9e95f2cc86ff9846098ed36c6ee5f478@192.168.5.70:55555",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	"enode://4d98470f306fa082cb526eff963a04aaa553b21eda2ede35ef5923a183efb8c28acfb06cb72a07ec6710ce239da4ce12c2d3eaf50b4adf463848d572859cc9b5@13.125.92.15:41171",
	"enode://cabe572bf008020cc125cd2a5353af57166adc5f14ee79261d9ca1f7c610a2a806e4b024ee17a609b95f09fc2bf26601433a5b0e156d5d18736004ae836ab550@34.237.211.223:41171",
	"enode://361616802dd85bede35d557fc53080eac9a637d7c2f76d7e7619ca7922f328ca3f3a5ee11a1522af68ad5bcf99e88398eb2566efdb0e201c5e1159b3a66aab5a@121.141.157.228:41172",
}
