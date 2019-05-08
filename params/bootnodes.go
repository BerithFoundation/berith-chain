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
	"enode://58e60ac860f35cf900b1fe92d51e419652651cff0cd7de085e5f415762a99b26c0c334d9a4a219ffd3f5032d42e5a0f446e1fca5b15a484f17442123f68d4aec@127.0.0.1:40404",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	"enode://58e60ac860f35cf900b1fe92d51e419652651cff0cd7de085e5f415762a99b26c0c334d9a4a219ffd3f5032d42e5a0f446e1fca5b15a484f17442123f68d4aec@127.0.0.1:40404",
	// "enode://58e60ac860f35cf900b1fe92d51e419652651cff0cd7de085e5f415762a99b26c0c334d9a4a219ffd3f5032d42e5a0f446e1fca5b15a484f17442123f68d4aec@121.141.157.230:40071",
	// "enode://e68ad1a7aea09aeb79992d00b6fe2a1790ebc37dcfee0b36a1b60362ac0f054d075ab0213b833b74f2acd5d471cadc4552f42cf9ea4e2098ab1ce87f94ad2b8f@121.141.157.228:40072",
	// "enode://0feec30cb54994d94a4c2306ac65f01de850a97a72889ffa8e4e93fe80699324960dd90408e5b53890ca3d75296a0165b73a665dce51da943759f9657866d528@34.237.211.223:40071",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	// "enode://58e60ac860f35cf900b1fe92d51e419652651cff0cd7de085e5f415762a99b26c0c334d9a4a219ffd3f5032d42e5a0f446e1fca5b15a484f17442123f68d4aec@121.141.157.230:40071",
	// "enode://e68ad1a7aea09aeb79992d00b6fe2a1790ebc37dcfee0b36a1b60362ac0f054d075ab0213b833b74f2acd5d471cadc4552f42cf9ea4e2098ab1ce87f94ad2b8f@121.141.157.228:40072",
	// "enode://0feec30cb54994d94a4c2306ac65f01de850a97a72889ffa8e4e93fe80699324960dd90408e5b53890ca3d75296a0165b73a665dce51da943759f9657866d528@34.237.211.223:40071",
}
