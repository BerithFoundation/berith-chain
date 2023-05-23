package types

import (
	"encoding/json"
	"testing"

	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/common/hexutil"
	"github.com/stretchr/testify/require"
)

func TestDecodeReceipt(t *testing.T) {
	fields := map[string]interface{}{
		"blockHash":         "0x6b803c4557260126b1a2532b138e842d0234f6aa7856ee13361a70cfae1f7fde",
		"blockNumber":       "0xbbb",
		"contractAddress":   "0x4d5950549ce92d462938d49185f14221352fe768",
		"cumulativeGasUsed": "0x36b2d9",
		"from":              "0x2345bf77d1de9eacf66fe81a09a86cfab212a542",
		"gasUsed":           "0x36b2d9",
		"logs": []*Log{
			{
				Address: common.HexToAddress("0x4d5950549ce92d462938d49185f14221352fe768"),
				Topics: []common.Hash{
					common.BytesToHash([]byte("0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d")),
					common.BytesToHash([]byte("0x0000000000000000000000000000000000000000000000000000000000000000")),
					common.BytesToHash([]byte("0x0000000000000000000000002345bf77d1de9eacf66fe81a09a86cfab212a542")),
					common.BytesToHash([]byte("0x0000000000000000000000002345bf77d1de9eacf66fe81a09a86cfab212a542")),
				},
				Data:        []byte{},
				BlockNumber: hexutil.MustDecodeUint64("0xbbb"),
				TxHash:      common.BytesToHash([]byte("0x5d7a888e2036a2f363edd143d984ba877f051c30945db74811de3efd657b8d1c")),
				TxIndex:     0,
				BlockHash:   common.BytesToHash([]byte("0x6b803c4557260126b1a2532b138e842d0234f6aa7856ee13361a70cfae1f7fde")),
				Index:       0,
				Removed:     false,
			},
		},
		"logsBloom":        "0x00000004000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000100000000000000000000000000000020000000000002000000800000000000000000000000000000000000000000000000000000200000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000100000000000020000000000000000000000000000000000000000000020000000002100000000000",
		"root":             "0x",
		"status":           "0x1",
		"to":               "",
		"transactionHash":  "0x5d7a888e2036a2f363edd143d984ba877f051c30945db74811de3efd657b8d1c",
		"transactionIndex": "0x0",
		"type":             "0x0",
	}
	input, _ := json.Marshal(fields)
	type Receipt struct {
		Type              *hexutil.Uint64 `json:"type,omitempty"`
		PostState         *hexutil.Bytes  `json:"root"`
		Status            *hexutil.Uint64 `json:"status"`
		CumulativeGasUsed *hexutil.Uint64 `json:"cumulativeGasUsed" gencodec:"required"`
		Bloom             *Bloom          `json:"logsBloom"         gencodec:"required"`
		Logs              []*Log          `json:"logs"              gencodec:"required"`
		TxHash            *common.Hash    `json:"transactionHash" gencodec:"required"`
		ContractAddress   *common.Address `json:"contractAddress"`
		GasUsed           *hexutil.Uint64 `json:"gasUsed" gencodec:"required"`
		BlockHash         *common.Hash    `json:"blockHash,omitempty"`
		BlockNumber       *hexutil.Big    `json:"blockNumber,omitempty"`
		TransactionIndex  *hexutil.Uint   `json:"transactionIndex"`
	}
	var dec Receipt
	if err := json.Unmarshal([]byte(input), &dec); err != nil {
		require.NoError(t, err)
	}
}
