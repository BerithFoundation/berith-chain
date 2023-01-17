package berithapi

import (
	"context"

	"github.com/BerithFoundation/berith-chain/core/types"
)

type TxPoolArgs interface {
	setDefaults(ctx context.Context, b Backend) error
	toTransaction() *types.Transaction
}
