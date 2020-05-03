package configuration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/configuration/types"
)

// NewHandler returns the handlers for the configuration module
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgUpdateConfig:
			return handleUpdateConfig(ctx, msg, k)
		case types.MsgDeleteLevelFee:
			return handleDeleteLevelFee(ctx, msg, k)
		case types.MsgUpsertDefaultFee:
			return handleUpsertDefaultFee(ctx, msg, k)
		case types.MsgUpsertLevelFee:
			return handleUpsertLevelFee(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "unknown request")
		}
	}
}

func handleUpdateConfig(ctx sdk.Context, msg types.MsgUpdateConfig, k Keeper) (*sdk.Result, error) {
	// check if operation is allowed
	allowed := k.OperationAllowed(ctx, msg)
	if !allowed {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "quorum for configuration update not reached")
	}
	// if allowed update configuration
	k.SetConfig(ctx, msg.NewConfiguration)
	// TODO emit event
	return &sdk.Result{}, nil
}

func handleDeleteLevelFee(ctx sdk.Context, msg types.MsgDeleteLevelFee, k Keeper) (*sdk.Result, error) {
	// check if operation is allowed
	allowed := k.OperationAllowed(ctx, msg)
	if !allowed {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "quorum for fee deletion not reached")
	}
	fees := k.GetFees(ctx)
	// not checking int overflow for 32bit machines because I suppose
	// our signers who are the owners are not trying to play themselves
	fees.DeleteLengthFee(msg, int(msg.Level.Int64()))
	// update fee
	k.SetFees(ctx, fees)
	// success TODO emit event?
	return &sdk.Result{}, nil
}

func handleUpsertDefaultFee(ctx sdk.Context, msg types.MsgUpsertDefaultFee, k Keeper) (*sdk.Result, error) {
	// check if operation is allowed
	allowed := k.OperationAllowed(ctx, msg)
	if !allowed {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "quorum for fee default fee upsertion not reached")
	}
	// get current fees
	fees := k.GetFees(ctx)
	// update fee
	fees.UpsertDefaultFees(msg, msg.Fee)
	// save in state
	k.SetFees(ctx, fees)
	// success TODO emit event?
	return &sdk.Result{}, nil
}

func handleUpsertLevelFee(ctx sdk.Context, msg types.MsgUpsertLevelFee, k Keeper) (*sdk.Result, error) {
	// check if operation is allowed
	allowed := k.OperationAllowed(ctx, msg)
	if !allowed {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "quorum for fee level fee upsertion not reached")
	}
	// get current fees
	fees := k.GetFees(ctx)
	// update level fee
	fees.UpsertLengthFees(msg, int(msg.Level.Int64()), msg.Fee)
	// save in state
	k.SetFees(ctx, fees)
	// success TODO emit event?
	return &sdk.Result{}, nil
}
