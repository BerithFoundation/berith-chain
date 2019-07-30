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
	"enode://3740dd6fda15ae822a5e69b1cbdb735e813761ecfda4f19eba04593ce91c125a7471ea3d189448733da5d2c649f078749423ad384782d14d3acd40c35f4e3d86@34.237.211.223:41171",
	"enode://8142df22d4ca164db41f0cb6a66f439332ec8d9a799dbd1a5bbc990ca585fbc283dc949edb8a0962942accb0c361a1171dc16dd46abefd17b7dc8ef19e6894e8@13.125.92.15:41171",
	"enode://0f146ca27cbe111f38ecf82b93a87718c37c2f4ed0e0b9d31b015d6e0253d12710221a63a68b298951638d61909f794e9da23272c08730d909fb666137524c1d@13.209.204.31:55555",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	"enode://4d98470f306fa082cb526eff963a04aaa553b21eda2ede35ef5923a183efb8c28acfb06cb72a07ec6710ce239da4ce12c2d3eaf50b4adf463848d572859cc9b5@13.125.92.15:41171",
	"enode://cabe572bf008020cc125cd2a5353af57166adc5f14ee79261d9ca1f7c610a2a806e4b024ee17a609b95f09fc2bf26601433a5b0e156d5d18736004ae836ab550@34.237.211.223:41171",
	"enode://361616802dd85bede35d557fc53080eac9a637d7c2f76d7e7619ca7922f328ca3f3a5ee11a1522af68ad5bcf99e88398eb2566efdb0e201c5e1159b3a66aab5a@121.141.157.228:41172",
}
