package common

import (
	"errors"

	"blockwatch.cc/tzgo/codec"
	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"github.com/alis-is/tezpay/constants/enums"
)

type OpExecutionContext struct {
	Op         *codec.Op
	Transactor TransactorEngine
	result     OpResult
}

func InitOpExecutionContext(op *codec.Op, transactor TransactorEngine) *OpExecutionContext {
	return &OpExecutionContext{
		Op:         op,
		result:     nil,
		Transactor: transactor,
	}
}

func (ctx *OpExecutionContext) GetOpHash() tezos.OpHash {
	if ctx.result == nil {
		return tezos.ZeroOpHash
	}
	return ctx.result.GetOpHash()
}

func (ctx *OpExecutionContext) Dispatch(opts *rpc.CallOptions) error {
	result, err := ctx.Transactor.Dispatch(ctx.Op, opts)
	if err != nil {
		return err
	}
	ctx.result = result
	return err
}

func (ctx *OpExecutionContext) WaitForApply() error {
	if ctx.result == nil {
		return errors.New("operation was not dispatched yet")
	}
	return ctx.result.WaitForApply()
}

type ITransferArgs interface {
	GetTxKind() enums.EPayoutTransactionKind
	GetFAContract() tezos.Address
	GetFATokenId() tezos.Z
	GetDestination() tezos.Address
	GetAmount() tezos.Z
}

func InjectTransferContents(op *codec.Op, source tezos.Address, p ITransferArgs) error {
	switch p.GetTxKind() {
	case enums.PAYOUT_TX_KIND_FA1_2:
		if p.GetFAContract().Equal(tezos.ZeroAddress) || p.GetFAContract().Equal(tezos.InvalidAddress) {
			return errors.New("invalid contract address")
		}
		args := contract.NewFA1TransferArgs().WithTransfer(source, p.GetDestination(), p.GetAmount()).
			WithDestination(p.GetFAContract())
		op.WithContents(args.Encode())
	case enums.PAYOUT_TX_KIND_FA2:
		if p.GetFAContract().Equal(tezos.ZeroAddress) || p.GetFAContract().Equal(tezos.InvalidAddress) {
			return errors.New("invalid contract address")
		}
		args := contract.NewFA2TransferArgs().WithTransfer(source, p.GetDestination(), p.GetFATokenId(), p.GetAmount()).
			WithDestination(p.GetFAContract())
		op.WithContents(args.Encode())
	default:
		op.WithTransfer(p.GetDestination(), p.GetAmount().Int64())
	}
	return nil
}
