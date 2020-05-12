// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Contains the metrics collected by the downloader.

package downloader

import (
	"berith-chain/metrics"
)

var (
	headerInMeter      = metrics.NewRegisteredMeter("berith/downloader/headers/in", nil)
	headerReqTimer     = metrics.NewRegisteredTimer("berith/downloader/headers/req", nil)
	headerDropMeter    = metrics.NewRegisteredMeter("berith/downloader/headers/drop", nil)
	headerTimeoutMeter = metrics.NewRegisteredMeter("berith/downloader/headers/timeout", nil)

	bodyInMeter      = metrics.NewRegisteredMeter("berith/downloader/bodies/in", nil)
	bodyReqTimer     = metrics.NewRegisteredTimer("berith/downloader/bodies/req", nil)
	bodyDropMeter    = metrics.NewRegisteredMeter("berith/downloader/bodies/drop", nil)
	bodyTimeoutMeter = metrics.NewRegisteredMeter("berith/downloader/bodies/timeout", nil)

	receiptInMeter      = metrics.NewRegisteredMeter("berith/downloader/receipts/in", nil)
	receiptReqTimer     = metrics.NewRegisteredTimer("berith/downloader/receipts/req", nil)
	receiptDropMeter    = metrics.NewRegisteredMeter("berith/downloader/receipts/drop", nil)
	receiptTimeoutMeter = metrics.NewRegisteredMeter("berith/downloader/receipts/timeout", nil)

	stateInMeter   = metrics.NewRegisteredMeter("berith/downloader/states/in", nil)
	stateDropMeter = metrics.NewRegisteredMeter("berith/downloader/states/drop", nil)
)
