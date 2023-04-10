package l1

import (
	"context"

	"github.com/ethereum-optimism/optimism/op-node/client"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum-optimism/optimism/op-node/sources"
	"github.com/ethereum-optimism/optimism/op-program/host/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

func NewFetchingL1(ctx context.Context, logger log.Logger, cfg *config.Config) (derive.L1Fetcher, error) {
	rpc, err := client.NewRPC(ctx, logger, cfg.L1URL)
	if err != nil {
		return nil, err
	}

	return sources.NewL1Client(rpc, logger, nil, sources.L1ClientDefaultConfig(cfg.Rollup, cfg.L1TrustRPC, cfg.L1RPCKind))
}

type Prefetcher struct {
	logger log.Logger
}

func NewPrefetcher(logger log.Logger) *Prefetcher {
	return &Prefetcher{logger: logger}
}

func (o *Prefetcher) BlockByHash(blockHash common.Hash) (*types.Block, error) {
	panic("implement me")
}